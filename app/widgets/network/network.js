'use strict';

angular.module('dashboard').controller('NetworkWidgetController', [
  function() {
  }
]);

angular.module('dashboard').factory('NetworkWidget', [
  'Widget',
  function(Widget) {
    return function() {
      var self = new Widget('network-widget', 'Network');

      self.data = {
        interface: null,
        received: null,
        transmitted: null,
      };

      self.identifier = function() {
        return self.type + "_" + self.configuration.interface;
      };

      return self;
    };
  }
]);

angular.module('dashboard').directive('networkWidget', [
  function() {
    return {
      replace: true,
      templateUrl: 'widgets/network/network.html',
    };
  }
]);
