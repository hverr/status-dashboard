'use strict';

describe('widget services', function() {
  beforeEach(module('dashboard'));

  describe('Widget', function() {
    it('should construct properly', inject(function(Widget) {
      var w = new Widget('my-directive', 'Widget Name');

      expect(w.directive).to.equal('my-directive');
      expect(w.name).to.equal('Widget Name');

      expect(w.height).to.equal(1);
      expect(w.width).to.equal(1);
      expect(w.row).to.equal(0);
      expect(w.col).to.equal(0);

      expect(w.type).to.equal(null);
      expect(w.clientIdentifier).to.equal(null);
      expect(w.client).to.equal(null);

      expect(w.data).to.equal(null);
    }));

    it('should update properly', inject(function(Widget) {
      var w = new Widget('my-directive', 'Widget Name');

      expect(w.data).to.equal(null);

      var d = { 'message' : 'Hello World!' };
      w.update(d);
      expect(w.data).to.equal(d);
    }));

    it('should configure properly', inject(function(Widget) {
      var w = new Widget('my-directive', 'Widget Name');

      expect(w.configuration).to.equal(null);
      var c = { 'key' : 'value' };
      w.configure(c);
      expect(w.configuration).to.equal(c);
    }));

    it('should use type as identifier', inject(function(Widget) {
      var w = new Widget('my-directive', 'Widget Name');

      expect(w.identifier()).to.equal(w.type);

      w.type = 'widget-type';
      expect(w.identifier()).to.equal(w.type);
    }));

    it('should return a watch value', inject(function(Widget) {
      var w = new Widget('my-directive', 'Widget Name');
      w.type = 'widget-type';
      w.configuration = { key: 'widget-configuration' };

      var watchValue = w.watchValue();
      expect(watchValue.type).to.equal(w.type);
      expect(watchValue.configuration).to.equal(w.configuration);
    }));
  });

  describe('widgetFactory', function() {
    var widgetFactory;

    beforeEach(inject(function(_widgetFactory_) {
      widgetFactory = _widgetFactory_;
    }));

    it('should return null for an unknown widget', function() {
      expect(widgetFactory('non-existing')).to.equal(null);
    });

    it('should create widgets', function() {
      expect(widgetFactory('load')).not.to.equal(null);
      expect(widgetFactory('uptime')).not.to.equal(null);
      expect(widgetFactory('meminfo')).not.to.equal(null);
      expect(widgetFactory('current_time')).not.to.equal(null);
      expect(widgetFactory('current_date')).not.to.equal(null);
      expect(widgetFactory('connections')).not.to.equal(null);
      expect(widgetFactory('network')).not.to.equal(null);
    });
  });

  describe('widgetsManager', function() {
    var widgetsManager;

    beforeEach(inject(function(_widgetsManager_) {
      widgetsManager = _widgetsManager_;
    }));

    describe('adding widgets', function() {
      it('should throw an exception for unknown widgets', function() {
        var fn = function() {
          widgetsManager.add('test-client', {
            type: 'unknown-widget-type',
          });
        };
        expect(fn).to.throw(/Unknown widget type/);
      });

      it('should set the properties of the widget', function() {
        var w = widgetsManager.add('client-name', {
          type: 'load',
        });
        expect(w.clientIdentifier).to.equal('client-name');
        expect(w.type).to.equal('load');
        expect(w.available).to.equal(false);
      });

      it('should configure the widget', function() {
        var c = {
          'interface' : 'lo',
        };
        var w = widgetsManager.add('client-name', {
          type: 'network',
          configuration: c,
        });
        expect(w.configuration).to.equal(c);
      });

      it('should leave client-less widgets available', function() {
        var w = widgetsManager.add("", {
          type: 'current_time',
        });
        expect(w.available).to.equal(true);
      });

      it('should properly load json widgets', function() {
        var json = [
          {
            client: 'webserver',
            type: 'network',
            width: 2,
            height: 1,
            row: 1,
            col: 1,
            configuration: {
              interface: 'lo',
            },
          },
          {
            client: "",
            type: 'current_time',
            width: 2,
            height: 2,
            row: 1,
            col: 3,
          },
        ];

        var result = widgetsManager.addFrom(json);
        expect(result.length).to.equal(2);
        expect(result[0].client).to.equal('webserver');
        expect(result[0].type).to.equal('network');
        expect(result[0].name).to.equal('Network (lo)');
        expect(result[0].width).to.equal(2);
        expect(result[0].height).to.equal(1);
        expect(result[0].row).to.equal(1);
        expect(result[0].col).to.equal(1);

        expect(result[1].client).to.equal("");
        expect(result[1].type).to.equal('current_time');
        expect(result[1].width).to.equal(2);
        expect(result[1].height).to.equal(2);
        expect(result[1].row).to.equal(1);
        expect(result[1].col).to.equal(3);
      });
    });
  });
});
