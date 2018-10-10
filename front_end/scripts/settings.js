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
                updateConfig(config, 'settings', config.settings)
            }
        }
    })
};
