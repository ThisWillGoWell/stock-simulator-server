var vm_topBar;

// var allViews = $(".view");
// var dashboardView = $("#dashboard--view");
// var businessView = $("#business--view");
// var stocksView = $("#stocks--view");
// var investorsView = $("#investors--view");
// var settingsView = $("#settings--view");
// var researchView = $("#research--view");
// var storeView = $("#store--view");
// var currentViewName = $("#current-view");

// function renderContent(route) {
//     switch (route) {
//         case "dashboard":
//         allViews.removeClass("active");
//         dashboardView.addClass("active");
//         currentViewName[0].innerHTML = "Dashboard";
//         break;

//         case "business":
//         allViews.removeClass("active");
//         businessView.addClass("active");
//         console.log(currentViewName);
//         currentViewName[0].innerHTML = "Business";
//         break;

//         case "stocks":
//         allViews.removeClass("active");
//         stocksView.addClass("active");
//         currentViewName[0].innerHTML = "Stocks";
//         break;

//         case "investors":
//         allViews.removeClass("active");
//         investorsView.addClass("active");
//         currentViewName[0].innerHTML = "Investors";
//         break;

//         case "research":
//         allViews.removeClass("active");
//         researchView.addClass("active");
//         currentViewName[0].innerHTML = "Research";
//         break;

//         case "settings":
//         allViews.removeClass("active");
//         settingsView.addClass("active");
//         currentViewName.text = "Settings";
//         break;

//         case "perks":
//         allViews.removeClass("active");
//         storeView.addClass("active");
//         currentViewName[0].innerHTML = "Store";
//         break;
//     }
// }

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
      goToSettings: function(event) {
        
        console.log("goToSettings");
        //renderContent("settings");
        
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
        return vm_dash_tab.currUserPortfolio.level;
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