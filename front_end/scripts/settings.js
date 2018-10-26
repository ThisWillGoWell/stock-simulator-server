var vm_settings;

function load_settings_tab() {
    vm_settings = new Vue({
        el: '#settings--view',
        methods: {
            toggleChangePercentSetting: function() {
                // Get config                 
                var config = vm_config.config;
                // Change value
                config.settings.changePercent = !config.settings.changePercent;
                //Send update
                updateConfig(config, 'settings', config.settings);
            },
            toggleSellAllButton: function() {
                var config = vm_config.config;
                console.log(config)
                config.settings.sellAll = !config.settings.sellAll;
                updateConfig(config, 'settings', config.settings);
            },
            changeDisplayName: function() {
                // Get entered display name
                let new_name = $("#change-display-name").val();
        
                // Creating message that changes the users display name
                let msg = {
                    set: "display_name",
                    value: new_name
                };
                
                let callback = function(msg) {
                    if (msg.msg.success) {
                        notifyTopBar("Hi, " + new_name + "!", GREEN, msg.msg.success);
                    } else {
                        notifyTopBar(msg.msg.error, RED, msg.msg.success);
                    }
                };
                // Send through WebSocket
                console.log(JSON.stringify(msg));
                doSend("set", msg, callback);
        
                // Reset display name
                $("#change-display-name").val("");
            },
        }
    })
};