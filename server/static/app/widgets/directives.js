'use strict';

angular.module('dashboard').directive('widget', [
  function() {
    return {
      templateUrl : 'widgets/widget.html',
    };
  }
]);

angular.module('dashboard').directive('widgetDynamicInfo', [
  '$compile',
  function($compile) {
    return {
      link: function(scope, element) {
        var d = scope.widget.directive;
        var template = '<div ' + d + '></div>';
        element.replaceWith($compile(template)(scope));
      },
    };
  }
]);
