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

      self.genericName = 'Network';
      self.data = {
        interface: null,
        received: null,
        transmitted: null,
      };

      self.configure = function(c) {
        self.configuration = c;

        self.interface = self.configuration.interface;
        self.name = self.genericName + ' (' + self.interface + ')';
      };

      self.identifier = function() {
        return self.type + "_" + self.interface;
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
