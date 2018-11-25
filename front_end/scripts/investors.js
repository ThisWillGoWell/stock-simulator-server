var vm_investors_tab;

function load_investors_tab() {
  // Vue for all investors tab data
  vm_investors_tab = new Vue({
    el: "#investors--view",
    data: {
      sortBy: "name",
      sortDesc: 1,
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
      toggleFavorite: function(uuid) {
        favoriteInvestor(uuid);
      },
      createGraph: function(portfolioUUID) {
        let location = "#investorGraph" + portfolioUUID;
        createPortfolioGraph(portfolioUUID, location);
      },
      openTransferModal: function(user) {
        transferModal.showModal = true;
        transferModal.recipient_uuid = user.uuid;
        transferModal.recipient_name = user.name;
        toggleTransferModal();
      },
      isFavoriteInvestor: function(uuid) {
        try {
          return (vm_config.config.fav.users.indexOf(uuid) > -1);
        } catch (err) {
          console.error(err);
          return false;
        }
      },
      makeInvestorGraph: function(investor) {
        $('button[data-route="research"]').click();
        $('#query-type-container .investors').click();
        graphUsers[0].selectize.addItem(investor.uuid);
      }
    },
    computed: {
      investors: function() {
        var investors = Object.values(vm_portfolios.portfolios);
        // List of all ledger items
        var ledgerItems = Object.values(vm_ledger.ledger);
        
        investors.map(function(d) {
          // Augment investor data
          d.name = vm_users.users[d.user_uuid].display_name;
          if (vm_store.items !== undefined) {
            if (d.level === 0) {
              d.title = "Novice";
            } else {
              d.title = vm_store.levels.filter(l => l.level == d.level)[0].title;
            }
          }
          // Get all stocks
          d.stocks = ledgerItems.filter(
            l => (l.portfolio_id === d.uuid) & (l.amount !== 0)
            ); // ledgers can have amount == 0, filter them out
            // Augment stock data
            d.stocks = d.stocks.map(function(d) {
              d.ticker_id = vm_stocks.stocks[d.stock_id].ticker_id;
              d.stock_name = vm_stocks.stocks[d.stock_id].name;
              d.current_price = vm_stocks.stocks[d.stock_id].current_price;
              d.value = d.current_price * d.amount;
              return d;
            });
  
            return d;
          });

          // Sort investors
          let byCol = this.sortBy;
          let direction = this.sortDesc;
          if (byCol === "favorites") {
            favs = vm_config.config.fav.users;
            investors = investors.sort(function(a, b) {
              if (favs.indexOf(a.uuid) === favs.indexOf(b.uuid)) {
                return 0;
              }
              if (favs.indexOf(a.uuid) > -1) {
                return direction;
              } else {
                return -direction;
              }
            })
            return investors;
          }

          investors = investors.sort(function(a, b) {
            if (typeof(a[byCol]) == "string") {
              if (a[byCol].toLowerCase() > b[byCol].toLowerCase()) {
                return -direction;
              }
              if (a[byCol].toLowerCase() < b[byCol].toLowerCase()) {
                return direction;
              }
              return 0;
            } else {
              if (a[byCol] > b[byCol]) {
                return -direction;
              }
              if (a[byCol] < b[byCol]) {
                return direction;
              }
              return 0;
            }
          });
          return investors;
        }
      },
    });
    
    // Set investor row clicking
    $("table").on("click", "tr.investors", function(event) {
      //var ticker_id = $(this).find('.stock-ticker-id').attr('tid');
      //console.log("TID: "+ticker_id);
      //var stock = Object.values(vm_users.stocks).filter(d => d.ticker_id === ticker_id)[0];
      // Set show modal to true
      //transferModal.showModal = true;
      //transferModal.investor_uuid = stock.uuid;
      //toggleTransferModal();
    });
}