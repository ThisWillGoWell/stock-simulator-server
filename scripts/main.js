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


	// $('.mobile-nav-btn').click(function() {
    
 //        $('#side-menu').addClass('open');
 //    });

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