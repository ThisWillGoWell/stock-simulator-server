var vm_config;

function createConfig(config) {

    if (jQuery.isEmptyObject(config)) {
        config = {
            fav: {
                stocks: [],
                users: [],
            },
            settings: {
                changePercent: false,
                sellAll: false,
                realValuesSetting: false,
                ticker: false,
                audioAlert: true,
            }
        };
    }
    vm_config = new Vue({
        data: {
            config: config
        }
    });
    console.log("------ CONFIG ------");
    console.log(vm_config.config);
    
    // Set checkboxes according to config settings
    checkSettingsBoxes();
}

// Method coming from stocks table favorite star
function favoriteStock(uuid) {
    
    var config = vm_config.config;

    // is stock in the list already
    if (config.fav.stocks.indexOf(uuid) > -1) {
        // remove the favorited stock 
        config.fav.stocks.splice(config.fav.stocks.indexOf(uuid), 1);
    
    } else {
        // Add new favorite
        config.fav.stocks.push(uuid);        

    }
    updateConfig(config, 'fav', config.fav);

};


// Method coming from investors table favorite star
function favoriteInvestor(uuid) {
    
    var config = vm_config.config;

    // is stock in the list already
    if (config.fav.users.indexOf(uuid) > -1) {
        // remove the favorited stock 
        config.fav.users.splice(config.fav.users.indexOf(uuid), 1);

    } else {
        // Add new favorite
        config.fav.users.push(uuid);

    }
    updateConfig(config, 'fav', config.fav);

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


// Getters
function getRealValuesSetting() {
    if (vm_config !== undefined) {
        return vm_config.config.settings.realValuesSetting;
    } else return false;
}

function getSellAllSetting() {
    if (vm_config !== undefined) {
        return vm_config.config.settings.sellAll;
    } else return false;
}

function getTickerSetting() {
    if (vm_config !== undefined) {
        return vm_config.config.settings.ticker;
    } else return false;
}

function getAudioAlertSetting() {
    if (vm_config !== undefined) {
        return vm_config.config.settings.audioAlert;
    } else return false;
}

function getConfigSetting(field) {
    if (vm_config !== undefined) {
        return vm_config.config.settings[field];
    } else return false;
}