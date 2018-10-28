var vm_topBar;


function load_topbar_vue() {
    
  // Vue for username top right
  vm_topBar = new Vue({
    el: "#top-bar--container",
    methods: {
      logout: function(event) {
        // delete cookie
        // Get saved data from sessionStorage
        console.log("logout");
        sessionStorage.removeItem("token");
        sessionStorage.removeItem("auth_obj");
        // send back to index
        window.location.href = "/login.html";
      },
    },
    computed: {
      userDisplayName: function() {
          var currUserUUID = sessionStorage.getItem("uuid");
          if (vm_users.users[currUserUUID] !== undefined) {
              return vm_users.users[currUserUUID].display_name;
          }
          return "";
      },
      userLevel: function() {
        let lvl = vm_dash_tab.currUserPortfolio.level;
        return " level " + lvl;
      },
      
    },
  });

  
    $(".account-settings-btn").click(function() {
        // console.log("clicked");
        $("#top-bar--container .account-settings-menu--container").toggleClass(
            "open"
        );
    });

    $("#account-settings-menu-close-btn").click(function() {
        $("#top-bar--container .account-settings-menu--container").toggleClass(
            "open"
        );
    });
}