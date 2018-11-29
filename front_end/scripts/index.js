// TODO break these out another file


/* Highest level Vue data object */
var vm_stocks = new Vue({
  data: {
    stocks: {}
  },
  watch: {
    stocks: function() {
      updateResearchStocks();
    },
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
  },
  computed: {
    currentUser: function() {
      // Get userUUID of the person that is logged in
      var currentUser = sessionStorage.getItem("uuid");
      if (currentUser !== null) {
        // Have they been added to the users object yet?
        if (vm_portfolios.portfolios != undefined) {
          return Object.values(vm_portfolios.portfolios).filter(d => d.user_uuid === currentUser)[0];
        } else {
          return "";
        }
      }
    }
  }
});

// var vm_items = new Vue({
//   data: {
//     items: {}
//   },
//   computed: {
//     userItems: function() {
//       return this.items.map(function(i) {
//         switch(i.config) {
//           case 'personal_broker':
//             i.include
//             break;
//         }
//         console.log(i)
//       })
//     }
//   }

// });

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
  },
  watch: {
    users: function() {
      updateResearchUsers();
    },
  }
});

var vm_effects = new Vue({
  data: {
    effects: {}
  }
})

var vm_recordBook = new Vue({
  data: {
    records: {},
  },
});

var vm_recordEntry = new Vue({
  data: {
    entries: {},
  }
});


registerRoute("connect", function(msg) {
  console.log("login recieved");

  if (msg.msg.success) {
    console.log(msg);
    sessionStorage.setItem("uuid", msg.msg.uuid);
    createConfig(msg.msg.config);
    
    setTimeout(function() {
      $('#loader--container').addClass("exit");
    }, 1500);
    
  } else {
    let err_msg = msg.msg.err;
    console.log(err_msg);
    console.log(msg);
    window.location.href = "/login.html";
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
      msg.msg.object.changePercent = 0;
      Vue.set(vm_stocks.stocks, msg.msg.uuid, msg.msg.object);
      break;

    case "ledger":
      Vue.set(vm_ledger.ledger, msg.msg.uuid, msg.msg.object);
      break;

    case "user":
      Vue.set(vm_users.users, msg.msg.uuid, msg.msg.object);
      break;

    case "item":
      Vue.set(vm_items.items, msg.msg.uuid, msg.msg.object);
      break;

    case "effect":
      Vue.set(vm_effects.effects, msg.msg.uuid, msg.msg.object);
      break;

    case "notification":
      Vue.set(vm_notify.notes, msg.msg.uuid, msg.msg.object);
      // If notification is not seen, notify user based on note type
      if (!msg.msg.object.seen) {
        // Execute notification type
        routeNote[msg.msg.object.type](msg.msg.object);
      }  
      break;

    case "record_book":
      Vue.set(vm_recordBook.records, msg.msg.uuid, msg.msg.object);
      break;

    case "record_entry":
      msg.msg.object.time = Date(msg.msg.object.time);
      Vue.set(vm_recordEntry.entries, msg.msg.uuid, msg.msg.object);
      break;
  }
});

registerRoute("delete", function(msg) {
  switch (msg.msg.type) {
    case 'effect':
      Vue.delete(vm_effects.effects, msg.msg.uuid);
      break;

    case 'item':
      Vue.delete(vm_items.items, msg.msg.uuid);
      break;
  }
})


registerRoute("alert", function(msg) {
  console.log(msg);
});


$(document).ready(function() {
  
  load_store_tab(); // store.js
  load_settings_tab(); // settings.js
  load_dashboard_tab(); // dashboard.js
  load_stocks_tab(); // stocks.js
  load_investors_tab(); // investors.js
  load_research_tab(); //research.js
  load_topbar_vue(); // topbar.js
  load_sidebar_vue(); // sidebar.js
  load_chat_vue(); // chat.js
  load_modal_vues(); // modal.js


  // setTimeout(function() {
  //   checkUsedItems(); // Display item perks that are in use 
  // }, 500);


  console.log("------ USER ITEMS ------")
  console.log(vm_items.items);
  console.log("------ USERS ------");
  console.log(vm_users.users);
  console.log("------ STOCKS ------");
  console.log(vm_stocks.stocks);
  console.log("------ LEDGER ------");
  console.log(vm_ledger.ledger);
  console.log("------ PORTFOLIOS ------");
  console.log(vm_portfolios.portfolios);
  console.log("------ EFFECTS ------");
  console.log(vm_effects.effects);
  console.log("------ NOTIFICATIONS ------");
  console.log(vm_notify.notes);
  console.log("------ RECORD BOOK ------");
  console.log(vm_recordBook.records);
  console.log("------ RECORD ENTRY ------");
  console.log(vm_recordEntry.entries);

  /* Vues that are used to display data */

  // Vue for sidebar navigation
  var vm_nav = new Vue({
    el: "#nav",
    methods: {
      nav: function(event) {
        let route = event.currentTarget.getAttribute("data-route");

        renderContent(route);
      }
    }
  });


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




  $("table").on("click", "tr td i.material-icons.star", function(event) {
    var ticker_id = $(this)
      .find(".stock-ticker-id")
      .attr("tid");

    console.log("TID: " + ticker_id + " has been favorited");

    //var stock = Object.values(vm_users.stocks).filter(d => d.ticker_id === ticker_id)[0];

    // Set show modal to true

    //transferModal.investor_uuid = stock.uuid;

    //vm_stocks_tab.toggleFavorite();
  });

  $("table").on("click", "tr td i.material-icons.chart", function(event) {
    var ticker_id = $(this)
      .find(".stock-ticker-id")
      .attr("tid");

    console.log("SHOW CHART FOR: " + ticker_id);

    //var stock = Object.values(vm_users.stocks).filter(d => d.ticker_id === ticker_id)[0];

    // Set show modal to true

    //transferModal.investor_uuid = stock.uuid;

    //vm_stocks_tab.toggleFavorite();
  });

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
    $(event.currentTarget)
      .find(".card.item")
      .toggleClass("hover");
    console.log("hover");
  });

  // $(".buy-item-btn.item-disabled").hover(function(event) {
  //   $(this)
  //     .parent(".card.item")
  //     .removeClass("hover");
  // });

  

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

        // Adding percent change amount
        vm_stocks.stocks[targetUUID].changePercent = Number(findPercentChange(targetChange, currPrice));

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

  var itemUpdate = function(msg) {
    var targetUUID = msg.msg.uuid;
    msg.msg.changes.forEach(function(changeObject) {
      // Variables needed to update the ledger item
      var targetField = changeObject.field;
      var targetChange = changeObject.value;

      // Update ledger item
      vm_items.items[targetUUID][targetField] = targetChange;
    });
  };

  var effectUpdate = function(msg) {
    var targetUUID = msg.msg.uuid;
    msg.msg.changes.forEach(function(changeObject) {
      var targetField = changeObject.field;
      var targetChange = changeObject.value;

      vm_effects.effects[targetUUID][targetField] = targetChange;
    });
  }

  var notificationUpdate = function(msg) {
    var targetUUID = msg.msg.uuid;
    msg.msg.changes.forEach(function(changeObject) {
      // Variables needed to update the ledger item
      var targetField = changeObject.field;
      var targetChange = changeObject.value;

      // Update ledger item
      vm_notify.notes[targetUUID][targetField] = targetChange;
    });
  }

  var recordBookUpdate = function(msg) {
    var targetUUID = msg.msg.uuid;
    msg.msg.changes.forEach(function(changeObject) {
      var targetField = changeObject.field;
      var targetChange = changeObject.value;

      vm_recordBook.records[targetUUID][targetField] = targetChange;
      console.log(vm_recordBook.records);
    })
  };


  registerRoute("update", function(msg) {
    var updateRouter = {
      stock: stockUpdate,
      ledger: ledgerUpdate,
      portfolio: portfolioUpdate,
      user: userUpdate,
      item: itemUpdate,
      effect: effectUpdate,
      notification: notificationUpdate,
      record_book: recordBookUpdate,
    };
    updateRouter[msg.msg.type](msg);
  });


  var allViews = $(".view");
  var dashboardView = $("#dashboard--view");
  var businessView = $("#business--view");
  var stocksView = $("#stocks--view");
  var investorsView = $("#investors--view");
  var settingsView = $("#settings--view");
  var researchView = $("#research--view");
  var storeView = $("#store--view");
  var currentViewName = $("#current-view");
  var currentRoute = "dashboard";

  function renderContent(route) {
      switch (route) {
          case "dashboard":
          if(route !== currentRoute) {
            currentRoute = route;
            allViews.removeClass("active");
            dashboardView.addClass("active");
            currentViewName[0].innerHTML = "Dashboard";
            TweenMax.from(currentViewName, 0.2, {ease: Back.easeOut.config(1.7), x:-10, opacity:0});
          } 
          break;

          case "business":
          if(route !== currentRoute) {
            currentRoute = route;
            allViews.removeClass("active");
            businessView.addClass("active");
            currentViewName[0].innerHTML = "Business";
            TweenMax.from(currentViewName, 0.2, {ease: Back.easeOut.config(1.7), x:-10, opacity:0});
          }
          break;

          case "stocks":
          if(route !== currentRoute) {
            currentRoute = route;
            allViews.removeClass("active");
            stocksView.addClass("active");
            currentViewName[0].innerHTML = "Stocks";
            TweenMax.from(currentViewName, 0.2, {ease: Back.easeOut.config(1.7), x:-10, opacity:0});
          }
          break;

          case "investors":
          if(route !== currentRoute) {
            currentRoute = route;
            allViews.removeClass("active");
            investorsView.addClass("active");
            currentViewName[0].innerHTML = "Investors";
            TweenMax.from(currentViewName, 0.2, {ease: Back.easeOut.config(1.7), x:-10, opacity:0});
          }
          break;

          case "research":
          if(route !== currentRoute) {
            currentRoute = route;
            allViews.removeClass("active");
            researchView.addClass("active");
            currentViewName[0].innerHTML = "Research";
            $('#research-graph-stock-select-selectized').focus();
            TweenMax.from(currentViewName, 0.2, {ease: Back.easeOut.config(1.7), x:-10, opacity:0});
          }
          break;

          case "settings":
          if(route !== currentRoute) {
            currentRoute = route;
            allViews.removeClass("active");
            settingsView.addClass("active");
            currentViewName[0].innerHTML = "Settings";
            TweenMax.from(currentViewName, 0.2, {ease: Back.easeOut.config(1.7), x:-10, opacity:0});
          }
          break;

          case "perks":
          if(route !== currentRoute) {
            currentRoute = route;
            allViews.removeClass("active");
            storeView.addClass("active");
            currentViewName[0].innerHTML = "Store";
            TweenMax.from(currentViewName, 0.2, {ease: Back.easeOut.config(1.7), x:-10, opacity:0});
          }
          break;
      }
  }

  init();
  
});
