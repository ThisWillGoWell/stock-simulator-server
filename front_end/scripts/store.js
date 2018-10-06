var vm_store;


function load_store_tab() {
    vm_store = new Vue({
        el: '#store--view',
        data: {

        },
        methods: {
            level_up: level_up,
        },
        computed: {
            currUserLevel: function() {
                console.log(vm_dash_tab.currUserPortfolio.level)
                return vm_dash_tab.currUserPortfolio.level;
            }
        }
    })
};



function level_up() {
    
    // Set callback
    REQUESTS[REQUEST_ID] = function (msg) {
        level_up_response(msg.msg.success, vm_store.currUserLevel);
    };

    // Send message
    doSend("level_up", {}, REQUEST_ID.toString());// REDO REQUEST ID CALC EVERYWHERE
    
    REQUEST_ID++;
}

function level_up_response(success, level) {
    if (success) {
        notify("Congrats you are level " + (Number(level) + 1), success);
    } else {
        notify("Error leveling up, not enough money.", success);
    }
}