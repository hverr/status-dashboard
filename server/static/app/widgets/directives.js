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
  '$log',
  function($compile, $log) {
    return {
      link: function(scope, element) {
        $log.debug('scope.widget:', scope.widget);
        var d = scope.widget.directive;
        var template = '<div ' + d + '></div>';
        element.replaceWith($compile(template)(scope));
      },
    };
  }
]);
