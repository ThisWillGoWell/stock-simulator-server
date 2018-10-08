var vm_store;

function load_store_tab() {
    vm_store = new Vue({
        el: '#store--view',
        data: {

        },
        methods: {
            level_up: level_up,
            purchaseItem: purchaseItem,
        },
        computed: {
            currUserLevel: function() {
                console.log(vm_dash_tab.currUserPortfolio.level)
                return vm_dash_tab.currUserPortfolio.level;
            }
        }
    })
};

function purchaseItem() {
    // Set callback
    var callback = function (msg) {
        if (msg.msg.o.success) {
            console.log("nothing for purchaseItem success callback");
            console.log(msg);
            
        } else {

            var message = msg.msg.o.err;
            var color = RED;

            notifyTopBar(message, color, msg.msg.o.success);
        }
    };

    var msg = {
        "action": "buy",
        "o": {
            "item_name": "insider"
        }
    };

    doSend("item", msg, callback);

};

function level_up() {
    
    // Set callback
    var callback = function (msg) {
        level_up_response(msg.msg, vm_store.currUserLevel);
    };

    // Send message
    doSend("level_up", {}, callback);// REDO REQUEST ID CALC EVERYWHERE
};

function level_up_response(msg, level) {
    if (msg.success) {
        var message = "You are now level " + (level + 1) + ".";
        notifyTopBar(message, GREEN)
    } else {
        var message = msg.err;

        notifyTopBar(message, RED, msg.success);
    }
};