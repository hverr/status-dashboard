'use strict';

describe('widget controllers', function() {
  beforeEach(module('dashboard'));

  describe('GridController', function() {
    describe('adding widgets', function() {
      it('should add widgets when requested', inject(function($controller, widgetsManager, widgetFactory) {
        var $scope = {
          $watch: function() {},
          widgets: [],
        };

        var $rootScope = {
          eventHandlers : {},
          $on: function(name, f) {
            $rootScope.eventHandlers[name] = f;
          },
        };

        widgetsManager = {
          addWidgetRequestEvent: widgetsManager.addWidgetRequestEvent,

          addRequests: [],
          add: function(client, widget) {
            expect(client).to.not.equal(null);
            expect(widget).to.not.equal(null);
            widgetsManager.addRequests.push([client, widget]);

            return widgetFactory(widget.type);
          },

          start: function() {},
          update: function() {},
        };

        var GridController = $controller('GridController', {
          $scope: $scope,
          $rootScope: $rootScope,
          widgetsManager: widgetsManager,
        });

        expect(GridController).to.not.equal(null);

        var handler = $rootScope.eventHandlers[widgetsManager.addWidgetRequestEvent];
        expect(handler).to.be.a('function');

        handler(widgetsManager.addWidgetRequestEvent, 'test-client', {type:'load'});
        expect(widgetsManager.addRequests.length).to.equal(1);
        expect($scope.widgets.length).to.equal(1);
        expect($scope.widgets[0].client).to.equal('test-client');
        expect($scope.widgets[0].col).to.equal(0);
        expect($scope.widgets[0].row).to.equal(0);
      }));
    });
  });
});
