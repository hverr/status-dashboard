'use strict';

describe('network widget', function() {

  beforeEach(module('dashboard'));
  beforeEach(module('testTemplates'));

  describe('NetworkWidgetController', function() {
    var NetworkWidgetController;

    beforeEach(inject(function($controller) {
      NetworkWidgetController = $controller('NetworkWidgetController');
    }));

    it('should construct properly', function() {
      expect(NetworkWidgetController).to.not.equal(null);
    });
  });

  describe('NetworkWidget', function() {
    var NetworkWidget;

    beforeEach(inject(function(_NetworkWidget_) {
      NetworkWidget = _NetworkWidget_;
    }));

    it('should construct properly', function() {
      var w = new NetworkWidget();

      expect(w.directive).to.equal('network-widget');
      expect(w.name).to.equal('Network');
      expect(w.data).to.have.property('interface');
      expect(w.data).to.have.property('received');
      expect(w.data).to.have.property('transmitted');
    });

    it('should configure properly', function() {
      var w = new NetworkWidget();
      var conf = {
        interface: 'lo',
      };

      w.configure(conf);
      expect(w.configuration).to.equal(conf);
      expect(w.interface).to.equal(conf.interface);
      expect(w.name).to.equal('Network (lo)');
    });

    it('should use the correct identifier', function() {
      var w = new NetworkWidget();
      var conf = {
        interface : 'lo',
      };
      w.configure(conf);
      w.type = 'network';

      expect(w.identifier()).to.equal('network_lo');
    });
  });

  describe('networkWidget', function() {
    var element, scope;

    beforeEach(inject(function($compile, $rootScope) {
      scope = $rootScope.$new();
      element = $compile('<div network-widget></div>')(scope);

      scope.widget = {
        data : {
          interface: 'lo',
          received: 2048,
          transmitted: 6,
        },
      };
      scope.data = scope.widget.data;
      scope.$digest();
    }));

    it('should bind properly', function() {
      expect(element.html()).to.contain('lo');
      expect(element.html()).to.contain('2.0 kB/s');
      expect(element.html()).to.contain('6.0 B/s');
    });
  });
});
