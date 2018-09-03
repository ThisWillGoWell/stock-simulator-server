$( document ).ready(function() {



	// Entrance Anims 
	$('#login-container').addClass('entrance-anim ');

	var scrollVal = 0;

	$( window ).on('resize', function(){
		screenWidth = $(window).width();
	    if (screenWidth > 1018) {
		
		} 
		
	});

	$(document).scroll(function() {

		scrollVal = $(document).scrollTop();

	    //console.log("SCROLL: "+scrollVal);


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

		// let attr = selector.attr('disable');

		// if (typeof attr !== typeof undefined && attr !== false) {
		// 	selector.prop("disabled", false); // Element(s) are now enabled.
		// } 
	}

	let input_login_uid = $('#login-uid');
	let input_login_pw = $('#login-pw');
	let input_login_submit = $('#input-login-submit');
	let input_create_uid = $('#create-uid');
	let input_create_pw = $('#create-pw');
	let input_create_submit = $('#input-create-submit');

	$('.create-account-btn').click(function() {
   

        input_login_uid.prop('disabled', true);
        input_login_pw.prop('disabled', true);
        input_login_submit.prop('disabled', true);

        input_create_uid.prop('disabled', false);
        input_create_pw.prop('disabled', false);
        input_create_submit.prop('disabled', false);

        console.log("clicked create account");

        // disableInput(input_login_uid);
        // disableInput(input_login_pw);

        // enableInput(input_create_uid);
        // enableInput(input_create_pw);

    });

 //    $('#menu-close-btn').click(function() {
 //    	$('#side-menu').removeClass('open');
        
 //    });

	// $("#publications-link").click(function() {
 //        $('html, body').animate({
 //            scrollTop: $("#publications").offset().top - 100
 //        }, 300);
 //    });

	/* particlesJS.load(@dom-id, @path-json, @callback (optional)); */
	// particlesJS.load('particles-js', '../assets/particles.json', function() {
	//   console.log('callback - particles.js config loaded');
	// });

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

});