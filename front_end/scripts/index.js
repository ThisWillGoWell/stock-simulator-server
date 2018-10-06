// TODO break these out another file


/* Highest level Vue data object */
var config = new Vue({
  data: {
    config: {}
  }
});

var vm_stocks = new Vue({
  data: {
    stocks: {}
  }
});

var vm_ledger = new Vue({
  data: {
    ledger: {}
  }
});

var vm_portfolios = new Vue({
  data: {
    portfolios: {}
  }
});

var vm_users = new Vue({
  data: {
    users: {},
    currentUser: auth_uuid
  },
  methods: {
    getCurrentUser: function() {
      // Get userUUID of the person that is logged in
      var currentUser = sessionStorage.getItem("uuid");
      if (currentUser !== null) {
        // Have they been added to the users object yet?
        if (vm_users.users[currentUser]) {
          return vm_users.users[currentUser].display_name;
        } else {
          return "";
        }
      }
    }
  }
});


registerRoute("connect", function(msg) {
  console.log("login recieved");

  if (!msg.msg.success) {
    let err_msg = msg.msg.err;
    console.log(err_msg);
    console.log(msg);
    window.location.href = "/login.html";
  } else {
    console.log(msg);
    sessionStorage.setItem("uuid", msg.msg.uuid);
    Vue.set(config.config, msg.msg.uuid, msg.msg.config);
  }
});

registerRoute("object", function(msg) {
  switch (msg.msg.type) {
    case "portfolio":
      //console.log(msg.msg.object)
      Vue.set(vm_portfolios.portfolios, msg.msg.uuid, msg.msg.object);
      break;

    case "stock":
      // Add variables for stocks for vue module initialization
      msg.msg.object.change = 0;
      Vue.set(vm_stocks.stocks, msg.msg.uuid, msg.msg.object);
      break;

    case "ledger":
      Vue.set(vm_ledger.ledger, msg.msg.uuid, msg.msg.object);
      break;

    case "user":
      Vue.set(vm_users.users, msg.msg.uuid, msg.msg.object);
      break;
  }
});


registerRoute("alert", function(msg) {
  console.log(msg);
});


$(document).ready(function() {
  load_dashboard_tab(); // dashboard.js
  load_investors_tab(); // investors.js
  load_stocks_tab(); // stocks.js
  load_store_tab(); // store.js
  load_topbar_vue(); // topbar.js
  load_notifications(); // notifications.js
  load_sidebar_vue(); // sidebar.js
  load_chat_vue(); // chat.js
  load_modal_vues(); // modal.js



  console.log("----- CONFIG -----");
  console.log(config.config);
  console.log("----- USERS -----");
  console.log(vm_users.users);
  console.log("------ STOCKS ------");
  console.log(vm_stocks.stocks);
  console.log("------ LEDGER ------");
  console.log(vm_ledger.ledger);
  console.log("------ PORTFOLIOS ------");
  console.log(vm_portfolios.portfolios);
  console.log("------ NOTIFICATIONS ------");
  console.log(vm_notify.notes);

  console.log(vm_topBar.userLevel)
  

  /* Vues that are used to display data */

  // Vue for sidebar navigation
  let vm_nav = new Vue({
    el: "#nav",
    methods: {
      nav: function(event) {
        let route = event.currentTarget.getAttribute("data-route");

        renderContent(route);
      }
    }
  });


  var notification_sound = new Audio();
  notification_sound.src = "assets/sfx_pling.wav";


  $(document).scroll(function() {
    scrollVal = $(document).scrollTop();
  });

  // $(".debug-title-bar button").click(function() {
  //   $("#debug-module--container").toggleClass("closed");
  //   //$('#debug-text-input').focus();
  // });


  // $(".debug-btn").click(function() {
  //   $("#debug-module--container").toggleClass("visible");
  // });




  // $("table").on("click", "tr td i.material-icons.star.unfilled", function(
  //   event
  // ) {
  //   var ticker_id = $(this)
  //     .find(".stock-ticker-id")
  //     .attr("tid");

  //   console.log("TID: " + ticker_id + " has been favorited");

  //   //var stock = Object.values(vm_users.stocks).filter(d => d.ticker_id === ticker_id)[0];

  //   // Set show modal to true

  //   //transferModal.investor_uuid = stock.uuid;

  //   vm_stocks_tab.toggleFavorite();
  // });

  $("thead tr th").click(function(event) {
    if (
      $(event.currentTarget)
        .find("i")
        .hasClass("shown")
    ) {
      $(event.currentTarget)
        .find("i")
        .toggleClass("flipped");
      // console.log("is asc");
    } else {
      $("thead tr th i").removeClass("shown");
      $(event.currentTarget)
        .find("i")
        .addClass("shown");
    }
  });

  // $(".buy-item-btn").click(function(event) {
  //   genericTextFieldModal.showModal = true;
  //   //transferModal.investor_uuid = stock.uuid;

  //   toggleGenericTextFieldModal();
  // });

  $(".buy-item-btn").hover(function(event) {
    $(this)
      .parent(".card.item")
      .toggleClass("hover");
  });

  

  $(document).keyup(function(e) {
    if (buySellModal.showModal === true) {
      if (e.keyCode === 27) {
        //toggleModal();
        buySellModal.closeModal();
      }
    }
  });


  var stockUpdate = function(msg) {
    var targetUUID = msg.msg.uuid;
    msg.msg.changes.forEach(function(changeObject) {
      // Variables needed to update the stocks
      var targetField = changeObject.field;
      var targetChange = changeObject.value;

      // if value to update is current price, calculate change
      if (targetField === "current_price") {
        // temp var for calculating price
        var currPrice = vm_stocks.stocks[targetUUID][targetField];
        // Adding change amount
        vm_stocks.stocks[targetUUID].change = targetChange - currPrice;
        // vm_stocks.stocks[targetUUID].change = Math.round((targetChange - currPrice) * 1000)/100000;

        // helper to color rows in the stock table
        var targetElem = $("tr[uuid=\x22" + targetUUID + "\x22]");
        var targetChangeElem = $(
          "tr[uuid=\x22" + targetUUID + "\x22] > td.stock-change"
        );

        if (targetChange - currPrice > 0) {
          targetChangeElem.removeClass("falling");
          targetChangeElem.addClass("rising");
        } else {
          targetChangeElem.removeClass("rising");
          targetChangeElem.addClass("falling");
        }
      }

      // Adding new current price
      vm_stocks.stocks[targetUUID][targetField] = targetChange;
    });
  };

  var ledgerUpdate = function(msg) {
    var targetUUID = msg.msg.uuid;
    msg.msg.changes.forEach(function(changeObject) {
      // Variables needed to update the ledger item
      var targetField = changeObject.field;
      var targetChange = changeObject.value;

      // Update ledger item
      vm_ledger.ledger[targetUUID][targetField] = targetChange;
    });
  };

  var portfolioUpdate = function(msg) {
    var targetUUID = msg.msg.uuid;
    msg.msg.changes.forEach(function(changeObject) {
      // Variables needed to update the ledger item
      var targetField = changeObject.field;
      var targetChange = changeObject.value;

      // Update ledger item
      vm_portfolios.portfolios[targetUUID][targetField] = targetChange;
    });
  };

  var userUpdate = function(msg) {
    var targetUUID = msg.msg.uuid;
    msg.msg.changes.forEach(function(changeObject) {
      // Variables needed to update the ledger item
      var targetField = changeObject.field;
      var targetChange = changeObject.value;

      // Update ledger item
      vm_users.users[targetUUID][targetField] = targetChange;
    });
  };

  
  registerRoute("update", function(msg) {
    var updateRouter = {
      stock: stockUpdate,
      ledger: ledgerUpdate,
      portfolio: portfolioUpdate,
      user: userUpdate
    };
    updateRouter[msg.msg.type](msg);
  });



var allViews = $(".view");
var dashboardView = $("#dashboard--view");
var businessView = $("#business--view");
var stocksView = $("#stocks--view");
var investorsView = $("#investors--view");
var futuresView = $("#futures--view");
var storeView = $("#store--view");
var currentViewName = $("#current-view");

function renderContent(route) {
    switch (route) {
        case "dashboard":
        allViews.removeClass("active");
        dashboardView.addClass("active");
        currentViewName[0].innerHTML = "Dashboard";
        break;

        case "business":
        allViews.removeClass("active");
        businessView.addClass("active");
        console.log(currentViewName);
        currentViewName[0].innerHTML = "Business";
        break;

        case "stocks":
        allViews.removeClass("active");
        stocksView.addClass("active");
        currentViewName[0].innerHTML = "Stocks";
        break;

        case "investors":
        allViews.removeClass("active");
        investorsView.addClass("active");
        currentViewName[0].innerHTML = "Investors";
        break;

        case "futures":
        allViews.removeClass("active");
        futuresView.addClass("active");
        currentViewName[0].innerHTML = "Futures";
        break;

        case "perks":
        allViews.removeClass("active");
        storeView.addClass("active");
        currentViewName[0].innerHTML = "Store";
        break;
    }
}

  // SOUND EFFECTS

  var notification_sound = new Audio();
  notification_sound.src = "assets/sfx_pling.wav";
  notification_sound.volume = 0.2;
});
