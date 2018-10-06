var vm_dash_tab;

function load_dashboard_tab() {
  // Vue for all dashboard data
  vm_dash_tab = new Vue({
    el: "#dashboard--view",
    data: {
      sortBy: "amount",
      sortDesc: 1
    },
    methods: {
      toPrice: formatPrice,
      // on column name clicks
      sortCol: function(col) {
        // If sorting by selected column
        if (vm_dash_tab.sortBy == col) {
          // Change sort direction
          // console.log(col);
          vm_dash_tab.sortDesc = -vm_dash_tab.sortDesc;
        } else {
          // Change sorted column
          vm_dash_tab.sortBy = col;
        }
        console.log(vm_dash_tab.sortBy);
        console.log(vm_dash_tab.currUserStocks);
      },
      createPortfolioGraph: function() {
        // Get curr user portfolioUUID
        let portfolioUUID = vm_dash_tab.currUserPortfolio.uuid;
        let location = "#portfolio-graph";
        createPortfolioGraph(portfolioUUID, location);
      }
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
              d.stock_roi =
                Number(d.stock_price) * Number(d.amount) -
                Number(d.investment_value);

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
            ownedStocks = ownedStocks.sort(function(a, b) {
              if (a[this.sortBy] > b[this.sortBy]) {
                return -this.sortDesc;
              }
              if (a[this.sortBy] < b[this.sortBy]) {
                return this.sortDesc;
              }
              return 0;
            });
            return ownedStocks;
          }
        }
        return [];
      }
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