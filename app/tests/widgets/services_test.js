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
    });
  });
});
