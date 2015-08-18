'use strict';

describe('load widget', function() {
  var LoadWidget;

  beforeEach(module('dashboard'));
  beforeEach(module('testTemplates'));


  beforeEach(inject(function(_LoadWidget_) {
    LoadWidget = _LoadWidget_;
  }));

  describe('LoadWidget', function() {
    it('should construct properly', function() {
      var w = new LoadWidget();

      expect(w.directive).to.equal('load-widget');
      expect(w.name).to.equal('Load');

      expect(w.data).to.have.property('cores');
      expect(w.data).to.have.property('one');
      expect(w.data).to.have.property('five');
      expect(w.data).to.have.property('fifteen');
    });
  });

  describe('LoadWidgetController', function() {
    var LoadWidgetController;
    var scope;

    beforeEach(inject(function($controller) {
      scope = {
        widget: new LoadWidget().data,
      };
      LoadWidgetController = $controller('LoadWidgetController', {
        $scope: scope
      });
    }));

    it('should detect large loads', function() {
      scope.widget.cores = 4;

      expect(scope.isLargeLoad(5.04)).to.equal(true);
      expect(scope.isLargeLoad("7.8")).to.equal(true);
    });

    it('should not detect normal loads', function() {
      scope.widget.cores = 4;

      expect(scope.isLargeLoad(2.95)).to.equal(false);
      expect(scope.isLargeLoad("1.06")).to.equal(false);
    });

    it('should handle invalid loads', function() {
      scope.widget.cores = 4;

      expect(scope.isLargeLoad("sheep")).to.equal(false);
      expect(scope.isLargeLoad(function() {})).to.equal(false);
      expect(scope.isLargeLoad(null)).to.equal(false);
    });
  });

  describe('loadWidget', function() {
    var element, scope;

    beforeEach(inject(function($compile, $rootScope) {
      scope = $rootScope.$new();
      element = $compile('<div load-widget></div>')(scope);

      scope.widget = {
        data : {
          cores: 4,
          one: "1.04",
          five: "5.03",
          fifteen: "2.03",
        },
      };
      scope.data = scope.widget.data;
      scope.$digest();

    }));

    it('should properly bind', function() {
      expect(element.html()).to.contain('4 cores');
      expect(element.html()).to.contain('1.04');
      expect(element.html()).to.contain('5.03');
      expect(element.html()).to.contain('2.03');
    });
  });
});
