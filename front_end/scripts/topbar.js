var topBar;

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

function load_topbar_vue() {
    
  // Vue for username top right
  topBar = new Vue({
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
      changeDisplayName: function() {
        // Get entered display name
        let new_name = $("#newDisplayName").val();

        // Creating message that changes the users display name
        let msg = {
          set: "display_name",
          value: new_name
        };

        REQUESTS[REQUEST_ID] = function(msg) {
          alert("Display_name changed to: " + new_name);
        };
        // Send through WebSocket
        console.log(JSON.stringify(msg));
        doSend("set", msg, REQUEST_ID.toString());

        REQUEST_ID++;

        // Reset display name
        $("#newDisplayName").val("");
      }
    },
    computed: {
      userDisplayName: function() {
        var currUserUUID = sessionStorage.getItem("uuid");
        if (vm_users.users[currUserUUID] !== undefined) {
          return vm_users.users[currUserUUID].display_name;
        }
        return "";
      }
    }
  });
}