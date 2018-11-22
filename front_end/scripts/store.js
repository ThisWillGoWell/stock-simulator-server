var vm_store;

function load_store_tab() {
    
    vm_store = new Vue({
        el: '#store--view',
        data: {
            store: {}
        },
        methods: {
            level_up: level_up,
            purchaseItem: purchaseItem,
            isNextLevel: function(lvl) {
                
                // if (this.currUserLevel >= lvl) {
                //     return true;
                // } else if (this.currUserLevel >= lvl+1){
                //     return true;
                // }
                return this.currUserLevel == (lvl-1);
            },
            isLowerThanLevel: function(lvl) {
                
                // if (this.currUserLevel >= lvl) {
                //     return true;
                // } else if (this.currUserLevel >= lvl+1){
                //     return true;
                // }
                return this.currUserLevel > lvl;
            },
            isHigherThanLevel: function(lvl) {
                
                // if (this.currUserLevel >= lvl) {
                //     return true;
                // } else if (this.currUserLevel >= lvl+1){
                //     return true;
                // }
                return this.currUserLevel < lvl;
            },
        },
        computed: {
            currUserLevel: function() {
                return vm_dash_tab.currUserPortfolio.level;
            },
            experience: function() {
                return this.store.exp;
            },
            items: function() {
                return this.store.abilities;
            }
        }
    });

    // Get item and level details
    var storeJSON = $.getJSON("../config/levels.json", function(data) {
        storeJSON = data;
        console.log(storeJSON)
    }).then(function(data) {
        console.log(data);
        Vue.set(vm_store.store, 'exp', data);
        // Vue.set(vm_store.store, 'abilities', data.abilities);
    });

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