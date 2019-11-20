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
                config.settings.sellAll = !config.settings.sellAll;
                updateConfig(config, 'settings', config.settings);
            },
            toggleRealValuesSetting: function() {
                var config = vm_config.config;
                config.settings.realValuesSetting = !config.settings.realValuesSetting;
                updateConfig(config, 'settings', config.settings);

                vm_dash_tab.realValueSetting = config.settings.realValuesSetting;
            },
            toggleTickerSetting: function() {
                var config = vm_config.config;
                config.settings.ticker = ! config.settings.ticker;
                updateConfig(config, 'settings', config.settings);
            },
            toggleAudioAlertSetting: function() {
                var config = vm_config.config;
                config.settings.audioAlert = ! config.settings.audioAlert;
                vm_chat.mute_notification_sfx = config.settings.audioAlert;
                console.log("CONFIG SOUND: "+ config.settings.audioAlert);
                updateConfig(config, 'settings', config.settings);
            },
            changeDisplayName: function() {
                // Get entered display name
                var new_name = $("#change-display-name").val();
                console.log(new_name)
                // Creating message that changes the users display name
                var msg = {
                    set: "display_name",
                    value: new_name
                };
                
                var callback = function(msg) {
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
            changePassword: function() {
                let new_pass = $("#change-password").val();
                let new_pass_2 = $("#change-password-confirm").val();

                if (new_pass === new_pass_2) {
                    let msg = {
                        set: 'password',
                        value: new_pass
                    }

                    let callback = function(msg) {
                        if (msg.msg.success) {
                            notifyTopBar("Password changed.", GREEN, msg.msg.success);
                        } else {
                            notifyTopBar(msg.msg.error, RED, msg.msg.success);
                        }
                    };

                    doSend("set", msg, callback);

                } else {
                    notifyTopBar("Passwords do not match.", RED)
                }
            },
        }
    })
};

// Set checkboxes according to config settings
function checkSettingsBoxes() {
    var settings = vm_config.config.settings;
    $('#percent-toggle-switch').prop('checked', settings.changePercent);
    $('#sell-all-toggle-switch').prop('checked', settings.sellAll);
    $('#actual-values-toggle-switch').prop('checked', settings.realValuesSetting);
    $('#audio-alert-toggle-switch').prop('checked', settings.audioAlert);
}