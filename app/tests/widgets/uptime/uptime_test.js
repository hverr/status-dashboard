'use strict';

describe('uptime widget', function() {
  var UptimeWidget;
  var $provide;

  beforeEach(module('dashboard'));
  beforeEach(module('testTemplates'));

  beforeEach(module(function(_$provide_) {
    $provide = _$provide_;
  }));

  beforeEach(inject(function(_UptimeWidget_) {
    UptimeWidget = _UptimeWidget_;
  }));

  describe('UptimeWidget', function() {
    it('should construct properly', function() {
      var w = new UptimeWidget();

      expect(w.directive).to.equal('uptime-widget');
      expect(w.name).to.equal('Uptime');

      expect(w.data).to.have.property('days');
      expect(w.data).to.have.property('hours');
      expect(w.data).to.have.property('minutes');
      expect(w.data).to.have.property('seconds');
    });

    it('should deep copy data', function() {
      var w = new UptimeWidget();

      var data = {
        days: 1,
        hours: 2,
        minutes: 3,
        seconds: 4,
      };

      w.update(data);
      w.data.hours = 0;

      expect(data.hours).to.equal(2);
    });
  });

  describe('uptimeWidget', function() {
    var $compile, $rootScope;

    beforeEach(inject(function(_$compile_, _$rootScope_)  {
      $compile = _$compile_;
      $rootScope = _$rootScope_;
    }));

    it('should properly bind', function() {
      $provide.value('oneSecondService', {
        add: function() {},
        remove: function() {},
      });

      var scope = $rootScope.$new();
      scope.data = { days: 1, hours: 2, minutes: 3, seconds: 4, };

      var w = $compile('<div uptime-widget></div>')(scope);
      scope.$digest();

      expect(w.html()).to.contain('1 day');
      expect(w.html()).to.contain('2h');
      expect(w.html()).to.contain('3m');
      expect(w.html()).to.contain('4s');
    });

    it('should pluralize days', function() {
      $provide.value('oneSecondService', {
        add: function() {},
        remove: function() {},
      });

      var scope = $rootScope.$new();
      scope.data = { days: 5, hours: 2, minutes: 3, seconds: 4, };

      var w = $compile('<div uptime-widget></div>')(scope);
      scope.$digest();

      expect(w.html()).to.contain('5 days');
    });

    it('should register for the one second service', function() {
      var didAdd = 0, didRemove = 0;

      $provide.value('oneSecondService', {
        add: function(f) {
          didAdd += 1;
          expect(f).to.be.a('function');
          return 1;
        },

        remove: function(handle) {
          didRemove += 1;
          expect(handle).to.equal(1);
        },
      });

      var scope = $rootScope.$new();
      scope.data = { days: 1, hours: 2, minutes: 3, seconds: 4, };

      var w = $compile('<div uptime-widget></div>')(scope);
      scope.$digest();

      expect(didAdd).to.equal(1);
      expect(didRemove).to.equal(0);

      w.find('div').first().trigger('$destroy');
      expect(didAdd).to.equal(1);
      expect(didRemove).to.equal(1);
    });

    it('should properly increase the time each second', function() {
      var increase;

      $provide.value('oneSecondService', {
        add: function(f) {
          increase = f;
        },
        remove: function() {},
      });

      var scope = $rootScope.$new();
      var w = $compile('<div uptime-widget></div>')(scope);
      scope.$digest();
      expect(increase).to.be.a('function');

      scope.data = { days: 1, hours: 2, minutes: 3, seconds: 54, };
      increase();
      scope.$digest();
      expect(scope.data.seconds).to.equal(55);
      expect(scope.data.minutes).to.equal(3);
      expect(scope.data.hours).to.equal(2);
      expect(scope.data.days).to.equal(1);
      expect(w.html()).to.contain('55s');

      scope.data = { days: 1, hours: 23, minutes: 59, seconds: 59, };
      increase();
      scope.$digest();
      expect(scope.data.seconds).to.equal(0);
      expect(scope.data.minutes).to.equal(0);
      expect(scope.data.hours).to.equal(0);
      expect(scope.data.days).to.equal(2);
      expect(w.html()).to.contain('2 days');
      expect(w.html()).to.contain('0h');
      expect(w.html()).to.contain('0m');
      expect(w.html()).to.contain('0s');
    });

    it('should handle increases with no data', function() {
      var increase;

      $provide.value('oneSecondService', {
        add: function(f) {
          increase = f;
        },
        remove: function() {},
      });

      var scope = $rootScope.$new();
      scope.data = null;
      $compile('<div uptime-widget></div>')(scope);
      scope.$digest();
      increase();
    });
  });
});
