var vm_config = new Vue({
    data: {
      config: {}
    }
});

{
    fav: {
        stocks: []
    }
}




// Method coming from stocks table favorite star
function favoriteStock(uuid) {
    console.log(uuid);
};


// Method coming from investors table favorite star
function favoriteInvestor(uuid) {
    console.log(uuid);
};


// Updating config
function updateConfig(new_config) {
    
    var msg = {
        set: 'config',
        value: new_config
    };

    var callback = function() {
        console.log("no callback for updateConfig");
    };

    doSend('set', msg, callback);
};