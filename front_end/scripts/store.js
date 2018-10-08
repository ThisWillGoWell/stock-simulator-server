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
        console.log("nothing for purchaseItem callback");
        console.log(msg)
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
        level_up_response(msg.msg.success, vm_store.currUserLevel);
    };

    // Send message
    doSend("level_up", {}, callback);// REDO REQUEST ID CALC EVERYWHERE
};

function level_up_response(success, level) {
    if (success) {
        notify("Congrats you are level " + (Number(level) + 1), success);
    } else {
        notify("Error leveling up, not enough money.", success);
    }
};