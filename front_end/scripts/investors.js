$(document).ready(function() {
  // Vue for all investors tab data
  var vm_investors_tab = new Vue({
    el: "#investors--view",
    methods: {
      toPrice: formatPrice,
      createGraph: function(portfolioUUID) {
        let location = "#investorGraph" + portfolioUUID;
        createPortfolioGraph(portfolioUUID, location);
      },
      openTransferModal: function(user) {
        transferModal.showModal = true;
        transferModal.recipient_uuid = user.uuid;
        transferModal.recipient_name = user.name;
        toggleTransferModal();
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
          // Get all stocks
          d.stocks = ledgerItems.filter(
            l => (l.portfolio_id === d.uuid) & (l.amount !== 0)
          ); // ledgers can have amount == 0
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
        return investors;
      }
    }
  });
});