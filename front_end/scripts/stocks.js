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
      toggleFavorite: function(uuid) {
        console.log(uuid);
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
          length: "100h"
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
            length: "100h"
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
            DrawLineGraph("#stock-graph", data);
          } else {
            setTimeout(drawGraphOnceDone, 100);
          }
        };

        setTimeout(drawGraphOnceDone, 100);
      }
    },
    computed: {
      sortedStocks: function() {
        if (Object.keys(vm_stocks.stocks).length !== 0) {
            // Turn to array and sort
            var stock_array = Object.values(vm_stocks.stocks);

            let byCol = this.sortBy;
            let direction = this.sortDesc;

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
          var mover = stocks.reduce((a, b) => (a.change > b.change ? a : b));
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
    $("#stock-list table").on("click", "tr.clickable", function(event) {
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
