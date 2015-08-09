'use strict';

angular.module('dashboard').directive('textWidget', [
  function() {
    return {
      link: function(scope, element) {
        element.addClass("text-content");
      },
    };
  }
]);
