$( document ).ready(function() {

    let url = "https://mockstarket.com";
    let querystring="";

    if( window.location.host.includes("localhost") ){
        if(window.location.port === "8080"){
            // force use localhost
            url = "http://localhost:8000";
            getToken("Will", "pass")
        }else if(window.location.port === "8081"){
            // force use dev
            url = "https://dev.mockstarket.com";
            querystring="?dev=1"; // this connects to dev instance
        } else if (window.location.port === "8082"){
            // force use prod
            url = "https://mockstarket.com"
        }
    } else if(window.location.host.includes("dev")){
        url = "https://dev.mockstarket.com";
        querystring="?dev=1";
    }


    let input_login_uid = $('#login-uid');
    let input_login_pw = $('#login-pw');
    let input_login_submit = $('#input-login-submit');
    let input_create_uid = $('#create-uid');
    let input_create_pw = $('#create-pw');
    let input_create_pw_confirm = $("#create-pw-confirm");
    let input_create_name = $('#create-name');
    let input_create_submit = $('#input-create-submit');
    var login_container = $("#login-container");

	console.log("login js");

    // Entrance Anims 
    $('#login-container').addClass('entrance-anim ');

    var scrollVal = 0;

    $( window ).on('resize', function(){
        screenWidth = $(window).width();
        if (screenWidth > 1018) {
        
        } 
        
    });

    
    function disableScroll() {
      if (window.addEventListener) // older FF
          window.addEventListener('DOMMouseScroll', preventDefault, false);
      window.onwheel = preventDefault; // modern standard
      window.onmousewheel = document.onmousewheel = preventDefault; // older browsers, IE
      window.ontouchmove  = preventDefault; // mobile
      document.onkeydown  = preventDefaultForScrollKeys;
    }

    function enableScroll() {
        if (window.removeEventListener)
            window.removeEventListener('DOMMouseScroll', preventDefault, false);
        window.onmousewheel = document.onmousewheel = null; 
        window.onwheel = null; 
        window.ontouchmove = null;  
        document.onkeydown = null;  
    }

    function disableWheelScroll(e){
        if(!e){ /* IE7, IE8, Chrome, Safari */ 
            e = window.event; 
        }
        if(e.preventDefault) { /* Chrome, Safari, Firefox */ 
            e.preventDefault(); 
        } 
        e.returnValue = false; /* IE7, IE8 */
    }

    var cardGroup = $('.card-group');

    cardGroup.bind('mousewheel DOMMouseScroll',function(){ 
        disableWheelScroll(); 
    });

    $('.create-user-btn ').click(function() {
        //enableScroll();
        $('.card-title-panel').addClass('show');

        var create_group = $("#create-user--container");
        
        console.log(cardGroup[0].scrollHeight);

        cardGroup.animate({
            scrollTop: cardGroup[0].scrollHeight
        }, 500, function(){
            $("input[type='text'][name='create-uid']").focus();
        });

        input_login_uid.prop('disabled', true);
        input_login_pw.prop('disabled', true);
        input_login_submit.prop('disabled', true);

        input_create_uid.prop('disabled', false);
        input_create_pw.prop('disabled', false);
        input_create_pw_confirm.prop('disabled', false);
        input_create_name.prop('disabled', false);
        input_create_submit.prop('disabled', false);
        

    });

    $('.card-title-panel i').click(function() {
        //enableScroll();
        $('.card-title-panel').removeClass('show');
        
        cardGroup.animate({
            scrollTop: 0
        }, 400, function(){
            $("input[type='text'][name='uid']").focus();
        });
        
        input_login_uid.prop('disabled', false);
        input_login_pw.prop('disabled', false);
        input_login_submit.prop('disabled', false);

        input_create_uid.prop('disabled', true);
        input_create_pw.prop('disabled', true);
        input_create_pw_confirm.prop('disabled', true);
        input_create_name.prop('disabled', true);
        input_create_submit.prop('disabled', true);

    });

    $('#mobile-consent--container .warning-content button').click(function() {
        $('#mobile-consent--container').remove();
    });

    function dismissWarning() {
        
    }

    function disableInput(selector) {

        let attr = selector.attr('disable');

        if (typeof attr !== typeof undefined && attr !== false) {
            $(selector).prop("disabled", false); // Element(s) are now enabled.
        } 
    }

    function enableInput(selector) {

        $(selector).prop('disabled', true);


    }

    

    $('.create-user-btn').click(function() {
   
        input_login_uid.prop('disabled', true);
        input_login_pw.prop('disabled', true);
        input_login_submit.prop('disabled', true);

        input_create_uid.prop('disabled', false);
        input_create_pw.prop('disabled', false);
        input_create_name.prop('disabled', false);
        input_create_submit.prop('disabled', false);

        console.log("clicked create user");

    });

 
    particlesJS("particles-js", {
      "particles": {
        "number": {
          "value": 63,
          "density": {
            "enable": true,
            "value_area": 1104.8066982851817
          }
        },
        "color": {
          "value": "#ffffff"
        },
        "shape": {
          "type": "circle",
          "stroke": {
            "width": 0,
            "color": "#000000"
          },
          "polygon": {
            "nb_sides": 5
          },
          "image": {
            "src": "img/github.svg",
            "width": 100,
            "height": 100
          }
        },
        "opacity": {
          "value": 0.5,
          "random": false,
          "anim": {
            "enable": false,
            "speed": 1,
            "opacity_min": 0.1,
            "sync": false
          }
        },
        "size": {
          "value": 3,
          "random": true,
          "anim": {
            "enable": false,
            "speed": 200,
            "size_min": 0.3,
            "sync": false
          }
        },
        "line_linked": {
          "enable": true,
          "distance": 150,
          "color": "#ffffff",
          "opacity": 0.4,
          "width": 1
        },
        "move": {
          "enable": true,
          "speed": 0.5,
          "direction": "none",
          "random": false,
          "straight": false,
          "out_mode": "out",
          "bounce": false,
          "attract": {
            "enable": false,
            "rotateX": 600,
            "rotateY": 1200
          }
        }
      },
      "interactivity": {
        "detect_on": "canvas",
        "events": {
          "onhover": {
            "enable": true,
            "mode": "repulse"
          },
          "onclick": {
            "enable": true,
            "mode": "push"
          },
          "resize": true
        },
        "modes": {
          "grab": {
            "distance": 400,
            "line_linked": {
              "opacity": 1
            }
          },
          "bubble": {
            "distance": 400,
            "size": 40,
            "duration": 2,
            "opacity": 8,
            "speed": 3
          },
          "repulse": {
            "distance": 200,
            "duration": 0.4
          },
          "push": {
            "particles_nb": 4
          },
          "remove": {
            "particles_nb": 2
          }
        }
      },
      "retina_detect": true
    });


    function authenticateUser(user, password) {
        var token = user + ":" + password;

        // Should i be encoding this value????? does it matter???
        // Base64 Encoding -> btoa
        var hash = btoa(token);

        return "Basic " + hash;
    };

    function getToken(user, password) {
        const Http = new XMLHttpRequest();
        Http.open("GET", url+"/api/token"+querystring , false);
        Http.setRequestHeader("Authorization", authenticateUser(user, password));
        Http.send();

        if (Http.status !== 200) {
            console.log(Http.responseText);
            loginFailed(Http.responseText);
            return null;
        } else {
            sessionStorage.setItem('token', Http.responseText);
            window.location.href = "/";
            return  Http.responseText;
        }
    };

    function createUser(user, password, nickname) {
        const Http = new XMLHttpRequest();
        Http.open("PUT", url+"/api/create"+querystring , false);
        Http.setRequestHeader("Authorization", authenticateUser(user, password));
        Http.setRequestHeader("DisplayName", nickname);
        Http.send();

        if (Http.status !== 200) {
            console.error(Http.responseText);
            return null;
        } else {
            sessionStorage.setItem('token', Http.responseText);
            window.location.href = "/";
            return  Http.responseText;
        }
    };

    function getLogin() {
        // Remove failed login message
        $('#login-failed-message').text("");

        let input_uid = $('#login-uid').val();
        let input_pw = $('#login-pw').val();

        getToken(input_uid, input_pw);
    }

    function loginFailed(msg) {
        $('#login-failed-message').text(msg);
        TweenMax.fromTo(login_container,0.10, {x:-20},{x:20,repeat:2,yoyo:true,ease:Sine.easeInOut,onComplete:function(){TweenMax.to(this.target,0.10,{x:0,ease:Elastic.easeOut})}})
    }

    $('#input-login-submit').click(function() {
        getLogin();
        console.log("login clicked");
    });

     $('#input-create-submit').click(function() {

        let input_uid = $('#create-uid').val().trim();
        let input_pw = $('#create-pw').val().trim();
        let input_create_pw_confirm = $('#create-pw-confirm').val().trim();
        let nickname = $('#create-name').val().trim();

        if(input_pw === input_create_pw_confirm) {
            createUser(input_uid, input_pw, nickname);
        }
        console.log("created user attempting login");
        
    });

    // On enter key try login
    $(document).keypress(function(e) {
        if(e.which == 13) {
            getLogin();
        }
    });


});