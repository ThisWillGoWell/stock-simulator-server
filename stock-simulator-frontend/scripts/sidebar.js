var sidebarCurrUser; 




function load_sidebar_vue() {
    
  // Vue for any sidebar data
  sidebarCurrUser = new Vue({
    el: "#stats--view",
    methods: {
      toPrice: formatPrice,
      realValuesSetting: getRealValuesSetting, 
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
      }
    }
  });
}