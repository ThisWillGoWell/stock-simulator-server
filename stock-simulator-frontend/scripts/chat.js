var vm_chat;


// SOUND EFFECTS
var notification_sound = new Audio();
notification_sound.src = "assets/sfx_pling.wav";
notification_sound.volume = 0.2;


var debug_feed = $("#debug-module--container .debug-message--list");

function appendNewMessage(msg, fromMe) {
    // if your chat is closed, add notification
    if (!vm_chat.showingChat) {
        vm_chat.unreadMessages = true;
        
        if (vm_config.config.settings.audioAlert) {
            console.log("SOUND: " + vm_config.config.settings.audioAlert);
            notification_sound.play();
        }
    }
    
    let chat_feed = $("#chat-module--container .chat-message--list");
    
    let msg_text = cleanInput(msg.body);
    let msg_author_display_name = msg.author_display_name;
    let msg_author_uuid = msg.author_uuid;
    let msg_timestamp = formatDate12Hour(new Date($.now()));

    let msg_template = "";
    let isMe = "";
    if (fromMe) {
        isMe = "is-me";
        msg_template =
        "<li " +
        msg_author_uuid +
        ">" +
        '				<div class="msg-timestamp">' +
        msg_timestamp +
        "</div>" +
        '				<div class="msg-username ' +
        isMe +
        '">' +
        msg_author_display_name +
        "</div>" +
        '				<div class="msg-text right">' +
        msg_text +
        "</div>" +
        "			</li>";
    } else {
        msg_template =
        "<li " +
        msg_author_uuid +
        ">" +
        '				<div class="msg-timestamp">' +
        msg_timestamp +
        "</div>" +
        '				<div class="msg-username ' +
        isMe +
        '">' +
        msg_author_display_name +
        "</div>" +
        '				<div class="msg-text">' +
        msg_text +
        "</div>" +
        "			</li>";
    }

    chat_feed.append(msg_template);
    chat_feed.animate(
        { scrollTop: chat_feed.prop("scrollHeight") },
        $("#chat-module--container .chat-message--list").height()
    );
}

function formatChatMessage(msg) {
    let timestamp = formatDate12Hour(new Date($.now()));
    // let message_body = $('#chat-module--container textarea').val();
    let message_body = msg.msg.message_body;
    let message_author = msg.msg.author;
    let isMe = false;

    if (vm_users.currentUser === message_author) {
        isMe = true;
        //console.log(isMe);
    } else {
        isMe = false;
    }

    let temp_msg = {
        author_uuid: message_author,
        author_display_name: vm_users.users[message_author].display_name,
        timestamp: timestamp,
        body: message_body
    };

    appendNewMessage(temp_msg, isMe);
}

  
$(document).keypress(function(e) {
    if ($("#chat-module--container textarea").val()) {
        if (e.which == 13) {
        let message_body = $("#chat-module--container textarea").val();

        var msg = {
            message_body: message_body
        };

        doSend("chat", msg);

        $("#chat-module--container textarea")
            .val()
            .replace(/\n/g, "");
        $("#chat-module--container textarea").val("");
        return false;
        }
    }
});

var cleanInput = input => {
    return $("<div/>")
        .text(input)
        .html();
};

function load_chat_vue() {

    registerRoute("chat", function(msg) {
        formatChatMessage(msg);
    });

    vm_chat = new Vue({
        el: "#chat-module--container",
        data: {
            showingChat: false,
            unreadMessages: false,
            mute_notification_sfx: true,
        },
        methods: {
            toggleChat: function() {
                this.showingChat = !this.showingChat;
                this.unreadMessages = false;
                $("#chat-module--container").toggleClass("closed");
                $("#chat-text-input").focus();
            },
            activeUsers: function() {
                // stop here later when not concating a string
                var online = Object.values(vm_users.users).filter(
                    d => d.active === true
                );

                var online_str = JSON.stringify(
                    online.map(d => d.display_name).join(", ")
                );
                return online_str.replace(/"/g, "");
            }
        },
        computed: {
            numActiveUsers: function() {
            return Object.values(vm_users.users).filter(d => d.active === true)
                .length;
            }
        },
        watch: {
            unreadMessages: function() {
                // make css changes here to show a notification for unread messages
                if (this.unreadMessages) {
                    console.log("unread messages");
                    $("#chat-module--container .chat-title-bar span").addClass("unread");
                    var link =
                    document.querySelector("link[rel*='icon']") ||
                    document.createElement("link");
                    link.type = "image/x-icon";
                    link.rel = "shortcut icon";
                    link.href = "assets/favicon_green_unread.png";
                    document.getElementsByTagName("head")[0].appendChild(link);
                    if (vm_config.config.settings.audioAlert) {
                        notification_sound.play();
                    }
                } else {
                    console.log("all messages read");
                    $("#chat-module--container .chat-title-bar span").removeClass(
                    "unread"
                    );
                    var link =
                    document.querySelector("link[rel*='icon']") ||
                    document.createElement("link");
                    link.type = "image/x-icon";
                    link.rel = "shortcut icon";
                    link.href = "assets/favicon_green.png";
                    document.getElementsByTagName("head")[0].appendChild(link);
                }
            },
            showingChat: function() {
                if (this.showingChat === false) {
                    $("#chat-text-input").blur();
                }
            },
        }
    });
    

}