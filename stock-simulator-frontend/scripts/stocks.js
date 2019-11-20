var vm_stocks_tab;
function load_stocks_tab() {

  // Vue for all stocks tab data
  vm_stocks_tab = new Vue({
    el: "#stocks--view",
    data: {
      sortBy: "ticker_id",
      sortDesc: 1,
      sortCols: ["ticker_id", "open_shares", "change", "current_price"],
      sortDirections: [-1, -1, -1, -1],
      reSort: 1
    },
    methods: {
      toPrice: formatPrice,
      formatPercent: formatPercent,
      toggleFavorite: function(uuid) {
        favoriteStock(uuid);
      },
      isFavoriteStock: function(uuid) {
        try {
          return (vm_config.config.fav.stocks.indexOf(uuid) > -1);
        } catch (err) {
          //console.error(err);
          return false;
        }
      },
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
    //   multiSort: function(col) {
    //     // if old first sort is the new first sort
    //     if (this.sortCols[0] === col) {
    //       // change sort direction
    //       this.sortDirections[0] *= -1;
    //     } else {
    //       // Where is the new sort column
    //       let ind = this.sortCols.indexOf(col);
    //       // Remove new column from old spot
    //       this.sortCols.splice(ind, 1);
    //       this.sortDirections.splice(ind, 1);
    //       // Push to the beginning of the array
    //       this.sortCols.unshift(col);
    //       this.sortDirections.unshift(1);
    //     }
    //     this.reSort++;
    //   },
      makeStockGraph: function(stock) {
        $('button[data-route="research"]').click();
        $('#query-type-container .stocks').click();
        graphStocks[0].selectize.addItem(stock.uuid);
      },
      amountOwned: function(stockUUID) {
        if (vm_dash_tab.currUserStocks !== undefined) {
          if (vm_dash_tab.currUserStocks.filter(d => d.stock_id === stockUUID).length === 1) {
            return vm_dash_tab.currUserStocks.filter(d => d.stock_id === stockUUID)[0].amount;
          } else {
            return 0;
          }
        }
        return 0;
      },
      createStockGraph: function(stockUUID) {
        let stock = Object.values(vm_stocks.stocks).filter(
          d => d.uuid === stockUUID
        )[0];
        console.log(stock);

        // Creating data object and adding tags
        let data = {
          data: {},
          tags: {
            title: stock.ticker_id,
            type: "stock"
          }
        };
        // Creating message
        let msg = {
          uuid: stockUUID,
          field: "current_price",
          num_points: 1000,
          use_cache: true,
          length: "6h"
        };

        // Store request on front end
        var callback = function(msg) {
          // Pull out the data and format it
          var points = msg.msg.points;
          points = points.map(function(d) {
            return { time: d[0], value: d[1] };
          });

          // Store the data
          data.data[stockUUID] = points;

          // Make note the data is available
          DrawLineGraph("#stock-list", data);
        };

        // Send message
        doSend("query", msg, callback);

      },
      createStocksGraph: function() {
        console.log("Creating stock graphs");
        // Store graphing data
        var data = {
          data: {},
          tags: {}
        };
        var responses = [];
        var requests = [];

        // Send data requests
        Object.keys(vm_stocks.stocks).forEach(function(stockUUID) {
          vm_stocks_tab.createStockGraph(stockUUID);
        });

        Object.keys(vm_stocks.stocks).forEach(function(stockUUID) {
          // Creating message
          let msg = {
            uuid: stockUUID,
            field: "current_price",
            use_cache: true,
            num_points: 1000,
            length: "6h"
          };

          // Store request on front end
          requests.push(REQUEST_ID.toString());
          var callback = function(msg) {
            // Pull out the data and format it
            var points = msg.msg.points;
            points = points.map(function(d) {
              return { time: d[0], value: d[1] };
            });

            // Store the data
            data.data[msg.msg.message.uuid] = points;

            // Make note the data is available
            responses.push(msg.request_id);
            // addToLineGraph('#portfolio-graph', points, field);
          };

          // Send message
          doSend("query", msg, callback);

        });

        var drawGraphOnceDone = null;

        var stillWaiting = true;

        drawGraphOnceDone = function() {
          if (requests.every(r => responses.indexOf(r) > -1)) {
            stillWaiting = false;
          }

          if (!stillWaiting) {
            console.log(data);
            DrawLineGraph("#stock-graph", data);
          } else {
            setTimeout(drawGraphOnceDone, 100);
          }
        };

        setTimeout(drawGraphOnceDone, 100);
      },
      openModal: function(ticker_id, buying) {

        console.log("TID: " + ticker_id);

        var stock = Object.values(vm_stocks.stocks).filter(
            d => d.ticker_id === ticker_id
        )[0];

        // Set show modal to true
        buySellModal.isBuying = buying;
        buySellModal.showModal = true;
        buySellModal.stock_uuid = stock.uuid;
        
        toggleModal();
      },
    //   buyOrder: function(tid, uuid) {
    //       console.log("BUY ORDER: "+uuid);
    //       // var ticker_id = $(this)
    //       //     .attr("tid");

    //       var ticker_id = tid;

    //       console.log("TID: " + ticker_id);

    //       var stock = Object.values(vm_stocks.stocks).filter(
    //           d => d.ticker_id === ticker_id
    //       )[0];
    //       buySellModal.setIsBuying(true);
    //       // Set show modal to true
    //       buySellModal.showModal = true;
    //       buySellModal.stock_uuid = uuid;

    //       toggleModal();
    //   },
    //   sellOrder: function(tid, uuid) {
    //     console.log("SELL ORDER: "+uuid);
    //     // var ticker_id = $(this)
    //     //     .attr("tid");

    //     var ticker_id = tid;

    //     console.log("TID: " + ticker_id);

    //     var stock = Object.values(vm_stocks.stocks).filter(
    //         d => d.ticker_id === ticker_id
    //     )[0];
        
    //     // Set show modal to true
    //     buySellModal.showModal = true;
    //     buySellModal.stock_uuid = uuid;

    //     toggleModal();
    //     $("#calc-btn-sell").addClass("fill");
    //     $("#calc-btn-buy").removeClass("fill");
    //     buySellModal.setIsBuying(false);
    // },
    },
    computed: {
      changePercentSetting: function() {
        if (vm_config === undefined) {
          return false;
        } else {
          return vm_config.config.settings.changePercent;
        }
      },
      sortedStocks: function() {
        if (Object.keys(vm_stocks.stocks).length !== 0) {
            
            let direction = this.sortDesc;
            
            // Turn to array and sort
            var stock_array = Object.values(vm_stocks.stocks);

            var byCol = this.sortBy;
            // Find which change are we sorting by
            if (byCol === "change" || byCol === "changePercent") {
              if (vm_config.config.settings.changePercent) {
                byCol = "changePercent";
              } else {
                byCol = "change";
              }
            } else if (byCol === "favorites") {
              favs = vm_config.config.fav.stocks;
              stock_array = stock_array.sort(function(a, b) {
                if (favs.indexOf(a.uuid) === favs.indexOf(b.uuid)) {
                  return 0;
                }
                if (favs.indexOf(a.uuid) > -1) {
                  return direction;
                } else {
                  return -direction;
                }
              })
              return stock_array;
            }

            // Sorting array
            stock_array = stock_array.sort(function(a, b) {
            if (a[byCol] > b[byCol]) {
                return -direction;
            }
            if (a[byCol] < b[byCol]) {
                return direction;
            }
            return 0;
          });
          return stock_array;
        }
        return [];
      },
    //   multiSortStocks: function() {
    //     if (Object.keys(vm_stocks.stocks).length !== 0) {
    //       function sorter(a, b, ind) {
    //         if (
    //           a[vm_stocks_tab.sortCols[ind]] > b[vm_stocks_tab.sortCols[ind]]
    //         ) {
    //           return vm_stocks_tab.sortDirections[ind];
    //         }
    //         if (
    //           a[vm_stocks_tab.sortCols[ind]] < b[vm_stocks_tab.sortCols[ind]]
    //         ) {
    //           return -vm_stocks_tab.sortDirections[ind];
    //         }
    //         if (ind === vm_stocks_tab.sortCols.length - 1) {
    //           return 0;
    //         } else {
    //           return sorter(a, b, ind + 1);
    //         }
    //       }

    //       // Get all stocks
    //       var stock_array = Object.values(vm_stocks.stocks);
    //       // Sort
    //       stock_array = stock_array.sort(function(a, b) {
    //         return sorter(a, b, 0);
    //       });

    //       return stock_array;
    //     }
    //     return [];
    //   },
      highestStock: function() {
        if (Object.values(vm_stocks.stocks).length === 0) {
          return "";
        } else {
          stocks = Object.values(vm_stocks.stocks).map(d => d);
          var highestStock = stocks.reduce(
            (a, b) => (a.current_price > b.current_price ? a : b)
          );
          
          return highestStock.ticker_id;
        }
      },
      mostChange: function() {
        if (Object.values(vm_stocks.stocks).length === 0) {
          return "";
        } else {
          stocks = Object.values(vm_stocks.stocks).map(d => d);
          if (vm_config.config.settings.changePercent) {
            var mover = stocks.reduce((a, b) => (a.changePercent > b.changePercent ? a : b));
          } else {
            var mover = stocks.reduce((a, b) => (a.change > b.change ? a : b));
          }
          // MAKE THIS HAPPEN ONLY IF A NEW STOCK TAKES OVER
          // TweenMax.to($(".stat-value i.most-change"), 0.2, {y: 15, ease:Bounce.easeOut});
          // TweenMax.to($(".stat-value i.most-change"), 0.2, {y: 0, delay:0.2});
          return mover.ticker_id;
        }
      },
      lowestStock: function() {
        if (Object.values(vm_stocks.stocks).length === 0) {
          return "";
        } else {
          stocks = Object.values(vm_stocks.stocks).map(d => d);
          var mover = stocks.reduce(
            (a, b) => (a.current_price < b.current_price ? a : b)
          );
          return mover.ticker_id;
        }
      }
    }
  });

    // Set stock row clicks
    // $("#stock-list table").on("click", "tr.clickable", function(event) {
        // var ticker_id = $(this)
        //     .find(".stock-ticker-id")
        //     .attr("tid");

        // console.log("TID: " + ticker_id);

        // var stock = Object.values(vm_stocks.stocks).filter(
        //     d => d.ticker_id === ticker_id
        // )[0];

        // // Set show modal to true
        // buySellModal.showModal = true;
        // buySellModal.stock_uuid = stock.uuid;

        // toggleModal();
    // });

    // $(".stat-value svg").on("mouseenter", function() {
    //   console.log("hovered");
    //   TweenMax.to(this, 0.3, {y: -12, ease:Bounce.easeOut})
    //   TweenMax.to(this, 0.2, {y: 0, delay:0.3})
    // })

    // $(".stat-value i").on("mouseenter", function() {
    //   console.log("hovered");
    //   TweenMax.to(this, 0.3, {y: -12, ease:Bounce.easeOut})
    //   TweenMax.to(this, 0.2, {y: 0, delay:0.3})
    // })
}
