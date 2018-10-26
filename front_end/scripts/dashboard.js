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
    },
    computed: {
      currUserPortfolio: function() {
        var currUserUUID = sessionStorage.getItem("uuid");
        if (vm_users.users[currUserUUID] !== undefined) {
          var currUserFolioUUID = vm_users.users[currUserUUID].portfolio_uuid;
          if (vm_portfolios.portfolios[currUserFolioUUID] !== undefined) {
            var folio = vm_portfolios.portfolios[currUserFolioUUID];
            folio.investments = folio.net_worth - folio.wallet;
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
                d.stock_roi = getROI(portfolio_uuid, d.stock_id, d.stock_price);

              // TODO: css changes done here talk to brennan about his \ux22 magic
              // helper to color rows in the stock table
              var targetChangeElem = $(
                "tr[uuid=\x22" +
                  d.stock_uuid +
                  "\x22].clickable > td.stock-change"
              );
              // targetChangeElem.addClass("rising");
              // if (d.stock_roi > 0) {
              // 	targetChangeElem.removeClass("falling");
              // 	targetChangeElem.addClass("rising");
              // } else if (d.stock_roi === 0) {
              // 	targetChangeElem.removeClass("falling");
              // 	targetChangeElem.removeClass("rising");
              // } else {
              // 	targetChangeElem.removeClass("rising");
              // 	targetChangeElem.addClass("falling");
              // }
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
      userItems: function() {
        var currUserUUID = sessionStorage.getItem("uuid");
        if (vm_users.users[currUserUUID] !== undefined) {
          var currUserFolioUUID = vm_users.users[currUserUUID].portfolio_uuid;
          var items = Object.values(vm_items.items).filter(d => d.portfolio_uuid === currUserFolioUUID);
          // Add used status
          console.log(items);
          items.map(function(d) {
            if (d.used) {
              d.used_status = 'Used';
            } else {
              d.used_status = 'Not Used';
            }
            return d;
          });

          return items;
        }
        return {};
      },

    }
  });

  // Set stock row clicks
  $("#owned-stocks").on("click", "tr.clickable", function(event) {
    var ticker_id = $(this)
        .find(".stock-ticker-id")
        .attr("tid");

    console.log("TID: " + ticker_id);

    var stock = Object.values(vm_stocks.stocks).filter(
        d => d.ticker_id === ticker_id
    )[0];

    // Set show modal to true
    buySellModal.showModal = true;
    buySellModal.stock_uuid = stock.uuid;

    toggleModal();
  });
}

function getROI(portfolio_uuid, stock_id, stock_price) {
  var userRecordsBooks = Object.values(vm_recordBook.records).filter(d => d.portfolio_uuid === portfolio_uuid);
  // Add stock id to record books 
  userRecordsBooks.forEach(function(d){
    d.stock_uuid = vm_ledger.ledger[d.ledger_uuid].stock_id;
    return d;
  });
  console.log(userRecordsBooks);
  var book = userRecordsBooks.filter(d => d.stock_uuid === stock_id)[0];

  var pricePaid = 0;
  var amountOwned = 0;

  book.buy_records.forEach(function(d) {
    amountOwned += d.AmountLeft;
    pricePaid += vm_recordEntry.entries[d.RecordUuid].result;
  });

  return amountOwned*stock_price + pricePaid;
}

function createPortfolioGraph(portfolioUUID, location) {
  // what it will be
  var uuids = [portfolioUUID, portfolioUUID];
  var fields = ['net_worth', 'wallet'];
  queryDrawGraph(location, uuids, fields);
}