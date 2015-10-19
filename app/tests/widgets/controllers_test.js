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

  describe('AddWidgetsDialogController', function() {
    describe('initializing', function() {
      var $scope;
      var AddWidgetsDialogController;

      beforeEach(inject(function($controller) {
        $scope = {};

        AddWidgetsDialogController = $controller('AddWidgetsDialogController', {
          $scope: $scope,
        });
      }));

      it('should properly initialize', function() {
        expect($scope.addWidget).to.be.a('function');
        expect($scope.refresh).to.be.a('function');
      });
    });

    describe('updating the contents', function() {
      var $scope;
      var $rootScope;
      var widgetsManager;
      var AddWidgetsDialogController;

      beforeEach(inject(function($controller, _widgetsManager_) {
        $scope = {};

        $rootScope = {
          eventHandlers: {},
          $on: function(name, f) {
            $rootScope.eventHandlers[name] = f;
          },
        };

        widgetsManager = _widgetsManager_;
        widgetsManager.availableClients = [
          {
            name : 'Web Server',
            identifier : 'webserver',
            availableWidgets : [
              {
                type : 'load',
                configuration : null,
              },
            ],
          }
        ];

        AddWidgetsDialogController = $controller('AddWidgetsDialogController', {
          $scope: $scope,
          $rootScope: $rootScope,
          widgetsManager : widgetsManager,
        });
      }));

      it('should update the contents when initialized', function() {
        expect($scope.clients.length).to.equal(2);
        expect($scope.clients[0].identifier).to.equal("");
        expect($scope.clients[1].identifier).to.equal('webserver');
      });

      it('should update the contents when available clients changed', function() {
        expect($rootScope.eventHandlers[widgetsManager.availableClientsChangedEvent]).to.be.a('function');

        widgetsManager.availableClients = [
          {
            name : 'FTP Server',
            identifier : 'ftp-server',
            availableWidgets : [
              {
                type : 'load',
                configuration : null,
              },
            ],
          }
        ];

        $rootScope.eventHandlers[widgetsManager.availableClientsChangedEvent]();
        expect($scope.clients.length).to.equal(2);
        expect($scope.clients[0].identifier).to.equal("");
        expect($scope.clients[1].identifier).to.equal('ftp-server');
      });

      it('should update when being refreshed', function() {
        var didUpdate = false;
        widgetsManager.updateAvailableClients = function() {
          didUpdate = true;
        };

        $scope.refresh();
        expect(didUpdate).to.equal(true);
      });
    });

    describe('adding widgets', function() {
      var AddWidgetsDialogController;
      var widgetsManager;
      var $scope;
      var $rootScope;

      beforeEach(inject(function($controller, _widgetsManager_) {
        widgetsManager = _widgetsManager_;
        $scope = {};
        $rootScope = {
          $on: function() {},

          emitted : {},
          $emit: function(name) {
            $rootScope.emitted[name] = arguments;
          },
        };
        AddWidgetsDialogController = $controller('AddWidgetsDialogController', {
          $scope : $scope,
          $rootScope : $rootScope,
        });
      }));

      it('should emit add widgets requests', function() {
        $scope.addWidget('webserver', {
          type: 'load'
        });

        expect(widgetsManager.addWidgetRequestEvent in $rootScope.emitted).to.equal(true);
      });
    });
  });
});
