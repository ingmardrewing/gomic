package fs

import (
	"fmt"

	"github.com/ingmardrewing/gomic/config"
)

type jsGen struct{}

func newJs() *jsGen {
	return &jsGen{}
}

func (j *jsGen) getJs() string {
	return fmt.Sprintf(js, cookiebar+newsletter)
}

func (j *jsGen) getAnalytics() string {
	if config.IsProd() {
		return analytics
	}
	return ""
}

func (j *jsGen) getDisqus(title string, disqusUrl string, disqusId string) string {
	return fmt.Sprintf(disqus_js, title, disqusUrl, disqusId)
}

var newsletter = `
$('.nl_container').staticFormHandler({
	  name: "devabodeNewsletterOffer",
	  fields: {
		"Email": {
		  "type":"input"
		}
	  },
	  url: "https://drewing.eu:16443/0.1/gomic/newsletter/add/",
	  display_condition: function(){ return $(window).scrollTop() > 200; },
});
`

var analytics = `
<script type="text/javascript">
//<![CDATA[
var gaProperty = 'UUA-49679648-1';
var disableStr = 'ga-disable-' + gaProperty;
if (document.cookie.indexOf(disableStr + '=true') > -1) {
  window[disableStr] = true;
}
function gaOptout() {
  document.cookie = disableStr + '=true; expires=Thu, 31 Dec 2099 23:59:59 UTC; path=/';
  window[disableStr] = true;
}

  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
  })(window,document,'script','//www.google-analytics.com/analytics.js','ga');

  ga('create', 'UA-49679648-1', 'devabo.de');
  ga('set', 'anonymizeIp', true);
  ga('require', 'displayfeatures');
  ga('require', 'linkid', 'linkid.js');
  ga('send', 'pageview');

//]]>
</script>`

var disqus_js = `
<div id="disqus_thread"></div>
<script>

var disqus_config = function () {
	this.page.title= "%s";
	this.page.url = 'https://DevAbo.de%s';
	this.page.identifier =  '%s';
};

(function() {
var d = document, s = d.createElement('script');
s.src = 'https://devabode.disqus.com/embed.js';
s.setAttribute('data-timestamp', +new Date());
(d.head || d.body).appendChild(s);
})();
</script>
<noscript>Please enable JavaScript to view the <a href="https://disqus.com/?ref_noscript">comments powered by Disqus.</a></noscript>
`

var js = `
(function($){
  $.fn.staticFormHandler = function(options){

    var defaults = {
        name: 'staticFormHandler',
        fields:{},
        url:"",
		intro_txt: "<h3>Want to get new pages via e-mail?</h3><p>Sign up to receive new DevAbo.de-pages via e-mail:<br>Just enter your e-mail address in the field below and click on &bdquo;Yes, sign me up&ldquo;.</p>",
        confirmation_txt:"<p>Thank you! You are almost ready - just click on the confirmation link sent to you via e-mail.</p>",
        error_txt:"Sorry, couldn't connect to the server. If you try again later, it might work.",
        display_condition:function(){ return false; },
        ask_only_once: true
      },
      plugin = this;

    this.opt = function(field){
      return options[field] || defaults[field];
    };

    this.createForm = function(){
      var $f = $("<form action="+ plugin.opt("action") + ">")
	  $f.append(plugin.opt('intro_txt'));
      $f.append(plugin.getFormFields());
	  $f.append(plugin.getDeclineButton());
	  $f.append(plugin.getSendButton());
      return $f;
    };

    this.readCookie = function(){
      c = document.cookie.split('; ');
      var i = c.length-1;
      while (i --> 0){
        var C = c[i].split('=');
        if(C[0] === plugin.opt('name')){
          return C[1];
        }
      }
    };

    this.getCookieExpireDate = function(){
      var now = new Date();
      var expDate = new Date();
      expDate.setYear(now.getFullYear()+20);
      return expDate.toString();
    };

    this.setCookie = function(){
      document.cookie = plugin.opt('name')+"=seen;expires=" + plugin.getCookieExpireDate();
    };

    this.alreadySeen = function(){
      return this.readCookie() === 'seen';
    };

    this.getInputField = function(f, c){
       return '<div class="nl_field"><label for="'+f+'"></label><input type="text" name="'+f+'" value="" id="'+f+'"></div>';
    };

    this.getSendButton = function(){
       var $btn = $('<a href="#" class="nl_button">Yes, sign me up!</a>');
      $btn.click(plugin.sendData);
       return $btn;
    };

    this.getDeclineButton = function(){
       var $btn = $('<a href="#" class="nl_button">No, thanks.</a>');
      $btn.click(plugin.close);
       return $btn;
    };

    this.getOkayCloseButton = function(){
       var $btn = $('<a href="#" class="nl_button">Okay</a>');
      $btn.click(plugin.close);
       return $btn;
    };

    this.getFormFields = function(){
      var fields = plugin.opt("fields"),
          fields_html = "";
      for( var f in fields){
        switch(fields[f].type){
          case "input":
            fields_html +=  plugin.getInputField(f, fields[f]);
            break;
        }
      }

      return $(fields_html);
    };

    this.clearInterval = function(){
      if( typeof(plugin.ti) !== 'undefined'){
        clearInterval(plugin.ti);
      }
    };

    this.getTriggerFunction = function(){
      return function(){
        if(plugin.opt('display_condition')()) {
          plugin.showForm();
          plugin.setCookie();
          plugin.clearInterval();
        }
      };
    };

    this.gatherData = function() {
      var fields = plugin.opt("fields"),
          data = {};
      for( var f in fields){
        data[f] = $("#" + f).val();
      }
      return JSON.stringify(data);
    };

    this.onAjaxError = function (err){
		plugin.each(function(){
			  $(this).find('.nl_error').remove();
			  $(this).find('.nl_text_container').prepend('<p class="nl_error">Please enter a valid e-mail address.</p>');
		  }) ;
    };

    this.onAjaxSuccess = function (data) {
	  if( data.Text === "address already registered"){
		plugin.each(function(){
			  $(this).find('.nl_error').remove();
			  $(this).find('.nl_text_container').prepend('<p class="nl_error">This e-mail address is already registered.</p>');
		  }) ;
	  }
	  else if( data.Text === "no address given"){
		plugin.each(function(){
			  $(this).find('.nl_error').remove();
			  $(this).find('.nl_text_container').prepend('<p class="nl_error">Please enter a valid e-mail address.</p>');
		  }) ;
	  }
	  else {
		  plugin.each(function(){
			  $c = $(this).find('.nl_text_container')
			  $c.children().remove();
			  $c.append( plugin.opt('confirmation_txt'));
			  $c.append(plugin.getOkayCloseButton());
		  }) ;
	  }
    };

    this.sendData = function (){
		var data =  plugin.gatherData();
       $.ajax({
        method: "PUT",
        url: plugin.opt('url'),
        data: data,
        dataType: "json",
        contentType: "application/json",
        error: plugin.onAjaxError,
        success: plugin.onAjaxSuccess
      });
      return false;
    };

    this.showForm = function(){
      plugin.each(function(){
        $(this).removeClass('nl_container_hidden');
		$c = $('<div class="nl_text_container">');
        $c.append( plugin.createForm());
		$(this).append($c);
      });
    };

    this.close= function(){
      plugin.each(function(){
        $(this).addClass('nl_container_hidden');
        $(this).children().remove();
      });
    };
  };
})(jQuery);

jQuery(document).ready(function() {
%s
});
`

var cookiebar = `
	function cli_show_cookiebar(p) {
		var Cookie = {
			set: function(name,value,days) {
				if (days) {
					var date = new Date();
					date.setTime(date.getTime()+(days*24*60*60*1000));
					var expires = "; expires="+date.toGMTString();
				}
				else var expires = "";
				document.cookie = name+"="+value+expires+"; path=/";
			},
			read: function(name) {
				var nameEQ = name + "=";
				var ca = document.cookie.split(';');
				for(var i=0;i < ca.length;i++) {
					var c = ca[i];
					while (c.charAt(0)==' ') {
						c = c.substring(1,c.length);
					}
					if (c.indexOf(nameEQ) === 0) {
						return c.substring(nameEQ.length,c.length);
					}
				}
				return null;
			},
			erase: function(name) {
				this.set(name,"",-1);
			},
			exists: function(name) {
				return (this.read(name) !== null);
			}
		};

		var ACCEPT_COOKIE_NAME = 'viewed_cookie_policy',
			ACCEPT_COOKIE_EXPIRE = 365,
			json_payload = p.settings;

		if (typeof JSON.parse !== "function") {
			console.log("CookieLawInfo requires JSON.parse but your browser doesn't support it");
			return;
		}
		var settings = JSON.parse(json_payload);

		var cached_header = jQuery(settings.notify_div_id),
			cached_showagain_tab = jQuery(settings.showagain_div_id),
			btn_accept = jQuery('#cookie_hdr_accept'),
			btn_decline = jQuery('#cookie_hdr_decline'),
			btn_moreinfo = jQuery('#cookie_hdr_moreinfo'),
			btn_settings = jQuery('#cookie_hdr_settings');

		cached_header.hide();
		if ( !settings.showagain_tab ) {
			cached_showagain_tab.hide();
		}

		var hdr_args = { };

		var showagain_args = { };
		cached_header.css( hdr_args );
		cached_showagain_tab.css( showagain_args );

		if (!Cookie.exists(ACCEPT_COOKIE_NAME)) {
			displayHeader();
		}
		else {
			cached_header.hide();
		}

		if ( settings.show_once_yn ) {
			setTimeout(close_header, settings.show_once);
		}
		function close_header() {
			Cookie.set(ACCEPT_COOKIE_NAME, 'yes', ACCEPT_COOKIE_EXPIRE);
			hideHeader();
		}

		var main_button = jQuery('.cli-plugin-main-button');
		main_button.css( 'color', settings.button_1_link_colour );

		if ( settings.button_1_as_button ) {
			main_button.css('background-color', settings.button_1_button_colour);

			main_button.hover(function() {
				jQuery(this).css('background-color', settings.button_1_button_hover);
			},
			function() {
				jQuery(this).css('background-color', settings.button_1_button_colour);
			});
		}
		var main_link = jQuery('.cli-plugin-main-link');
		main_link.css( 'color', settings.button_2_link_colour );

		if ( settings.button_2_as_button ) {
			main_link.css('background-color', settings.button_2_button_colour);

			main_link.hover(function() {
				jQuery(this).css('background-color', settings.button_2_button_hover);
			},
			function() {
				jQuery(this).css('background-color', settings.button_2_button_colour);
			});
		}

		cached_showagain_tab.click(function(e) {
			e.preventDefault();
			cached_showagain_tab.slideUp(settings.animate_speed_hide, function slideShow() {
				cached_header.slideDown(settings.animate_speed_show);
			});
		});

		jQuery("#cookielawinfo-cookie-delete").click(function() {
			Cookie.erase(ACCEPT_COOKIE_NAME);
			return false;
		});
		jQuery("#cookie_action_close_header").click(function(e) {
			e.preventDefault();
			accept_close();
		});

		function accept_close() {
			Cookie.set(ACCEPT_COOKIE_NAME, 'yes', ACCEPT_COOKIE_EXPIRE);

			if (settings.notify_animate_hide) {
				cached_header.slideUp(settings.animate_speed_hide);
			}
			else {
				cached_header.hide();
			}
			cached_showagain_tab.slideDown(settings.animate_speed_show);
			return false;
		}

		function closeOnScroll() {
			if (window.pageYOffset > 100 && !Cookie.read(ACCEPT_COOKIE_NAME)) {
				accept_close();
				if (settings.scroll_close_reload === true) {
					location.reload();
				}
				window.removeEventListener("scroll", closeOnScroll, false);
			}
		}
		if (settings.scroll_close === true) {
			window.addEventListener("scroll", closeOnScroll, false);
		}

		function displayHeader() {
			if (settings.notify_animate_show) {
				cached_header.slideDown(settings.animate_speed_show);
			}
			else {
				cached_header.show();
			}
			cached_showagain_tab.hide();
		}
		function hideHeader() {
			if (settings.notify_animate_show) {
				cached_showagain_tab.slideDown(settings.animate_speed_show);
			}
			else {
				cached_showagain_tab.show();
			}
			cached_header.slideUp(settings.animate_speed_show);
		}
	};

	function l1hs(str){if(str.charAt(0)=="#"){str=str.substring(1,str.length);}else{return "#"+str;}return l1hs(str);}

cli_show_cookiebar({
					settings: '{"animate_speed_hide":"500","animate_speed_show":"500","background":"#fff","border":"#444","border_on":true,"button_1_button_colour":"#000","button_1_button_hover":"#000000","button_1_link_colour":"#fff","button_1_as_button":true,"button_2_button_colour":"#333","button_2_button_hover":"#292929","button_2_link_colour":"#444","button_2_as_button":false,"font_family":"inherit","header_fix":false,"notify_animate_hide":true,"notify_animate_show":false,"notify_div_id":"#cookie-law-info-bar","notify_position_horizontal":"right","notify_position_vertical":"bottom","scroll_close":false,"scroll_close_reload":false,"showagain_tab":false,"showagain_background":"#fff","showagain_border":"#000","showagain_div_id":"#cookie-law-info-again","showagain_x_position":"100px","text":"#000","show_once_yn":false,"show_once":"10000"}'
});

`
