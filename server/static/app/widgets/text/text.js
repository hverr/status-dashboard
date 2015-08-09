'use strict';

angular.module('dashboard').directive('textWidget', [
  '$compile',
  function($compile) {
    return {
      link: function(scope, element) {
        var content = element.html();
        var replacement = '<div class="text-content"><div class="text">' + content + '</div></div>';
        element.replaceWith($compile(replacement)(scope));
      },
    };
  }
]);
