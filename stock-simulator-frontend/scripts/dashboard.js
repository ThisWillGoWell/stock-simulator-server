var vm_dash_tab;

function load_dashboard_tab() {
  // Vue for all dashboard data
  vm_dash_tab = new Vue({
    el: "#dashboard--view",
    data: {
      sortBy: "amount",
      sortDesc: 1,
      insiderStocks: [],
    },
    methods: {
      toPrice: formatPrice,
      // on column name clicks
      sortCol: function(col) {
        // If sorting by selected column
        if (this.sortBy == col) {
          // Change sort direction
          // console.log(col);
          this.sortDesc = -this.sortDesc;
        } else {
          // Change sorted column
          this.sortBy = col;
        }
      },
      createPortfolioGraph: function() {
        // Get curr user portfolioUUID
        let portfolioUUID = vm_dash_tab.currUserPortfolio.uuid;
        let location = "#portfolio-graph";
        // what will the graph be
        var uuids = [portfolioUUID, portfolioUUID];
        var fields = ['net_worth', 'wallet'];
        queryDrawGraph(location, uuids, fields);
      },
      useItem: function(item_uuid) {
        useItem(item_uuid);
      },
      sellAll: function(stock_id) {
        var amt = this.currUserStocks.filter(d => d.uuid === stock_id)[0].amount;
        var id = this.currUserStocks.filter(d => d.uuid === stock_id)[0].stock_id;
        sendTrade(id, (-1)*amt);
      },
      sellAllSetting: getSellAllSetting,
      realValuesSetting: getRealValuesSetting,
      buyOrder: function(tid, uuid) {
        console.log("BUY ORDER: "+uuid);
        var ticker_id = tid;
            

        console.log("TID: " + tid);

        var stock = Object.values(vm_stocks.stocks).filter(
            d => d.ticker_id === ticker_id
        )[0];

        // Set show modal to true
        buySellModal.showModal = true;
        buySellModal.stock_uuid = stock.uuid;

        toggleModal();
      },
      sellOrder: function(tid, uuid) {
          console.log("SELL ORDER: "+uuid);
          // var ticker_id = $(this)
          //     .attr("tid");

          var ticker_id = tid;

          console.log("TID: " + tid);

          var stock = Object.values(vm_stocks.stocks).filter(
              d => d.ticker_id === ticker_id
          )[0];
          
          // Set show modal to true
          buySellModal.showModal = true;
          buySellModal.stock_uuid = stock.uuid;

          toggleModal();
          
      },
    },
    computed: {
      currUserPortfolio: function() {
        var currUserUUID = sessionStorage.getItem("uuid");
        if (vm_users.users[currUserUUID] !== undefined) {
          var currUserFolioUUID = vm_users.users[currUserUUID].portfolio_uuid;
          if (vm_portfolios.portfolios[currUserFolioUUID] !== undefined) {
            var folio = vm_portfolios.portfolios[currUserFolioUUID];
            folio.investments = folio.net_worth - folio.wallet;

            // Adding real networth
            if (folio.stocksValue === undefined) folio.stocksValue = 0;
            if (folio.stocks === undefined) folio.stocks = [];
            prospectStockValues(folio.stocks, folio.stocksValue);

            return folio;
          }
        }
        return {};
      },
      currUserStocks: function() {
        var currUserUUID = sessionStorage.getItem("uuid");
        if (vm_users.users[currUserUUID] !== undefined) {
          // Current users portfolio uuid
          var portfolio_uuid = vm_users.users[currUserUUID].portfolio_uuid;

          
          // If objects are in ledger
          if (Object.keys(vm_ledger.ledger).length !== 0) {
            var ownedStocks = Object.values(vm_ledger.ledger).filter(
              d => d.portfolio_id === portfolio_uuid
            );
            
            // Remove stocks that user owns 0 of
            ownedStocks = ownedStocks.filter(d => d.amount !== 0);
            // Augmenting owned stocks
            ownedStocks = ownedStocks.map(function(d) {
              
              d.stock_ticker = vm_stocks.stocks[d.stock_id].ticker_id;
              d.stock_price = vm_stocks.stocks[d.stock_id].current_price;
              d.stock_value = Number(d.stock_price) * Number(d.amount);
              try {
                d.stock_roi = getROI(portfolio_uuid, d.stock_id, d.stock_price);
              }
              catch(err) {
                d.stock_roi = 0;
              }

              // TODO: css changes done here talk to brennan about his \ux22 magic
              // helper to color rows in the stock table
              var targetChangeElem = $(
                'tr[uuid="dash' + d.stock_uuid + '"].clickable > td.stock-change'
              );
              // targetChangeElem.addClass("rising");
              if (d.stock_roi > 0) {
              	targetChangeElem.removeClass("falling");
              	targetChangeElem.addClass("rising");
              // } else if (d.stock_roi === 0) {
              // 	targetChangeElem.removeClass("falling");
              // 	targetChangeElem.removeClass("rising");
              } else {
              	targetChangeElem.removeClass("rising");
              	targetChangeElem.addClass("falling");
              }
              return d;
            });

            // Sorting array of owned stocks
            
            let byCol = this.sortBy;
            let direction = this.sortDesc;

            ownedStocks = ownedStocks.sort(function(a, b) {
              if (a[byCol] > b[byCol]) {
                return -direction;
              }
              if (a[byCol] < b[byCol]) {
                return direction;
              }
              return 0;
            });

            return ownedStocks;
          }
        }
        return [];
      },
      userEffects: function() {
        var currUserUUID = vm_users.currentUser;
        if (vm_users.users[currUserUUID] !== undefined) {
          var currUserFolioUUID = vm_users.users[currUserUUID].portfolio_uuid;
          console.log(vm_effects.effects)
          var effects = Object.values(vm_effects.effects).filter(d => d.portfolio_uuid === currUserFolioUUID);
          console.log(effects);

          effects.forEach(function(e) {
            switch(e.title) {
              case "Personal Broker":
                e.desc = [
                  ["Purchase Fee", "$0"],
                  ["Sale Fee", "$0"]
                ]
                break;
              case "Base Effect":
                e.desc = [
                  ["Purchase Fee", "$" + formatPrice(e.buy_fee_amount)],
                  ["Sale Fee", "$" + formatPrice(e.sell_fee_amount)],
                  ["Profit Multiplier", (e.profit_multiplier * 100) + "%"],
                  ["Tax Rate of Profits", (e.tax_percent * 100) + "%"]
                ];
            }
          })

          return effects;
        }
        return {}; 
      },
      userItems: function() {
        var currUserUUID = sessionStorage.getItem("uuid");
        if (vm_users.users[currUserUUID] !== undefined) {
          var currUserFolioUUID = vm_users.users[currUserUUID].portfolio_uuid;
          var items = Object.values(vm_items.items).filter(d => d.portfolio_uuid === currUserFolioUUID);
          console.log(items)
          var ret = items.map(function(i) {
            var item = {
              name: i.name,
              config_id: i.config,
              uuid: i.uuid,
              portfolio_uuid: i.portfolio_uuid,
            };
            switch(i.config) {
              case 'personal_broker':
                item.duration = prettifyItemDuration(i.duration);
                item.desc = [
                  ['Purchase Fee', '$0'],
                  ['Sale Fee', '$0']
                ]
                break;
            }
            console.log(item)
            return item;
          })
          console.log(ret)
          return ret;
        }
        return {};
      },
      // sellAllSetting: function() {
      //   if (vm_config === undefined) {
      //     return false;
      //   } else {
      //     return vm_config.config.settings.sellAll;
      //   }
      // },
    }
  });

  // Set stock row clicks
  // $("#owned-stocks").on("click", "tr.clickable", function(event) {
  //   var ticker_id = $(this)
  //       .find(".stock-ticker-id")
  //       .attr("tid");

  //   console.log("TID: " + ticker_id);

  //   var stock = Object.values(vm_stocks.stocks).filter(
  //       d => d.ticker_id === ticker_id
  //   )[0];

  //   // Set show modal to true
  //   buySellModal.showModal = true;
  //   buySellModal.stock_uuid = stock.uuid;

  //   toggleModal();

  // });

  // Set stock row clicks
  // $("#owned-stocks").on("click", "tr.clickable", function(event) {
  //   var ticker_id = $(this)
  //       .find(".stock-ticker-id")
  //       .attr("tid");

  //   console.log("TID: " + ticker_id);

  //   var stock = Object.values(vm_stocks.stocks).filter(
  //       d => d.ticker_id === ticker_id
  //   )[0];

  //   // Set show modal to true
  //   buySellModal.showModal = true;
  //   buySellModal.stock_uuid = stock.uuid;

  //   toggleModal();

  // });

}







function getROI(portfolio_uuid, stock_id, stock_price) {
  var userRecordsBooks = Object.values(vm_recordBook.records).filter(d => d.portfolio_uuid === portfolio_uuid);
  // Add stock id to record books 
  userRecordsBooks.forEach(function(d) {
    d.stock_uuid = vm_ledger.ledger[d.ledger_uuid].stock_id;
    return d;
  });
  var book = userRecordsBooks.filter(d => d.stock_uuid === stock_id)[0];
  if (book !== undefined) {
    var pricePaid = 0;
    var amountOwned = 0;
  
    book.buy_records.forEach(function(d) {
      // Wait until entry has arrived
      if (vm_recordEntry.entries[d.RecordUuid] !== undefined) {
        pricePaid += vm_recordEntry.entries[d.RecordUuid].result;
      }
      amountOwned += d.AmountLeft;
    });
  
    return amountOwned*stock_price + pricePaid;
  } else return 0;
}

function createPortfolioGraph(portfolioUUID, location) {
  // what it will be
  var uuids = [portfolioUUID, portfolioUUID];
  var fields = ['net_worth', 'wallet'];
  queryDrawGraph(location, uuids, fields);
}



