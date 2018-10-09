var vm_config = new Vue({
    data: {
      config: {}
    }
});


// Method coming from stocks table favorite star
function favoriteStock(uuid) {
    console.log(uuid);
    
    var config = vm_config.config;
    // Add fav area if not included already
    if (config.fav === undefined) {
        config.fav = {};
    }
    // Add fav.stocks if not included already
    if (config.fav.stocks === undefined) {
        config.fav.stocks = [];
    }
    // larger than 5?
    if (config.fav.stocks.length > 10) {
        config.fav.stocks.length.pop();
    }
    // Add new favorite
    config.fav.stocks.unshift(uuid);
    config.fav.stocks
    
    updateConfig(config, 'stocks', config.fav.stocks);

};


// Method coming from investors table favorite star
function favoriteInvestor(uuid) {
    console.log(uuid);
    
    var config = vm_config.config;
    // Add fav area if not included already
    if (config.fav === undefined) {
        config.fav = {};
    }
    // Add fav.stocks if not included already
    if (config.fav.users === undefined) {
        config.fav.users = [];
    }
    // larger than 5?
    if (config.fav.users.length > 10) {
        config.fav.users.length.pop();
    }
    // Add new favorite
    config.fav.users.unshift(uuid);

    updateConfig(config, 'users', config.fav.users);

};


// Updating config
function updateConfig(new_config, new_key, new_value) {
    
    // Setting vue to create reactivity
    Vue.set(vm_config.config, new_key, new_value);

    var msg = {
        set: 'config',
        value: new_config
    };

    var callback = function(msg) {
        console.log(msg);
        if (msg.msg.success) {
            notifyTopBar("Success!", GREEN, msg.msg.success);
        } else {
            notifyTopBar("Uh oh!", RED, msg.msg.success);
        }
    };

    doSend('set', msg, callback);
};