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
    },
    computed: {
      investors: function() {
        var investors = Object.values(vm_portfolios.portfolios);
        // List of all ledger items
        var ledgerItems = Object.values(vm_ledger.ledger);
        
        investors.map(function(d) {
          // Augment investor data
          d.name = vm_users.users[d.user_uuid].display_name;
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

          investors = investors.sort(function(a, b) {
            if (a[byCol] > b[byCol]) {
              return -direction;
            }
            if (a[byCol] < b[byCol]) {
              return direction;
            }
            return 0;
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