$( document ).ready(function() {

    let token_path = "https://mockstarket.com/token";
    let input_login_uid = $('#login-uid');
    let input_login_pw = $('#login-pw');
    let input_login_submit = $('#input-login-submit');
    let input_create_uid = $('#create-uid');
    let input_create_pw = $('#create-pw');
    let input_create_pw_confirm = $("#create-pw-confirm");
    let input_create_name = $('#create-name');
    let input_create_submit = $('#input-create-submit');

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

    $('.create-account-btn ').click(function() {
        //enableScroll();
        $('.card-title-panel').addClass('show');

        var create_group = $("#create-account--container");
        
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

    function disableInput(selector) {

        let attr = selector.attr('disable');

        if (typeof attr !== typeof undefined && attr !== false) {
            $(selector).prop("disabled", false); // Element(s) are now enabled.
        } 
    }

    function enableInput(selector) {

        $(selector).prop('disabled', true);


    }

    

    $('.create-account-btn').click(function() {
   

        input_login_uid.prop('disabled', true);
        input_login_pw.prop('disabled', true);
        input_login_submit.prop('disabled', true);

        input_create_uid.prop('disabled', false);
        input_create_pw.prop('disabled', false);
        input_create_name.prop('disabled', false);
        input_create_submit.prop('disabled', false);

        console.log("clicked create account");

        

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
          "value": 5,
          "random": true,
          "anim": {
            "enable": false,
            "speed": 40,
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
          "speed": 2,
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
        Http.open("GET", token_path, false);
        Http.setRequestHeader("Authorization", authenticateUser(user, password));
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

    function attemptLogin() {
    	
    	let input_uid = $('#login-uid').val();
		let input_pw = $('#login-pw').val();

        getToken(input_uid, input_pw);

        // TODO delete if working
  //   	if(input_uid != '' && input_pw != '') {
  //           input_uid_trimmed = input_uid.trim();
  //           input_pw_trimmed = input_pw.trim();
  //   		auth_msg = {
  //   					action: "login",
		// 		        msg: {
		// 		        	"username": input_uid_trimmed, 
		// 		        	"password": input_pw_trimmed
		// 		        }
		// 		    };
  //           uid = input_uid_trimmed;
  //   	}

  //   	console.log(auth_msg);

  //   	try {
    		
	 //  		doSend(JSON.stringify(auth_msg));

  //           //doSend('{"action": "login", "msg": {"username": "Will", "password":"pass"}}');
		// }

		// catch(error) {
	 //  		console.error(error);
		//   	$('.login-err').text("Username or password is incorrect");
		// }
		
    }

    $('#input-login-submit').click(function() {
    	attemptLogin();
        console.log("login clicked");
    });

     $('#input-create-submit').click(function() {

        let input_uid = $('#create-uid').val();
        let input_pw = $('#create-pw').val();
        let input_create_pw_confirm = $('#create-pw-confirm').val();
        let input_name = $('#create-name').val();
        let create_auth_msg = {};

        if(input_uid != '' && input_pw != '') {
            input_uid_trimmed = input_uid.trim();
            input_pw_trimmed = input_pw.trim();

            input_create_pw_trimmed = input_create_pw_confirm.trim();
            input_name_trimmed = input_name.trim();

            if(input_pw_trimmed === input_create_pw_trimmed) {
                create_auth_msg = {
                    action: "new_account",
                    msg: {
                        "username": input_uid_trimmed, 
                        "password": input_pw_trimmed,
                        "display_name": input_name_trimmed
                    }
                };
            }
            

            uid = input_uid_trimmed;
        }

        console.log(create_auth_msg);

        try {
            
            doSend(JSON.stringify(create_auth_msg));

        }

        catch(error) {
            console.error(error);
            $('.login-err').text("Username or password is incorrect");
        }

        console.log("created account attempting login");
        
    });

    $(document).keypress(function(e) {
        
        if(e.which == 13) {
            attemptLogin();

        }
    
    });


});