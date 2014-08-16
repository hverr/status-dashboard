/** System **/
QUnit.test("System initializes correctly", function(assert) {
	var system = new System("System Name", "http://example.org");
	assert.ok(system != null);
	assert.equal(system.name, "System Name");
	assert.equal(system.updateURL, "http://example.org");
	assert.ok(system.components instanceof Array);
});

/** MemoryComponent **/
QUnit.test("Initialize MemoryComponent without json", function(assert) {
	var comp = new MemoryComponent();
	assert.ok(comp, "Component must be defined");
	assert.equal(comp.free, null);
	assert.equal(comp.used, null);
});

QUnit.test("Initialize MemoryComponent with wrong json", function(assert) {
	var invalid = {};
	assert.throws(function() { new MemoryComponent(invalid); });
});

QUnit.test("Initialize MemoryComponent with valid json", function(assert) {
	var valid = {
		memory : {
			total:500,
			used:200
		}
	};
	var comp = new MemoryComponent(valid);
	assert.ok(comp);
	assert.equal(comp.total, 500);
	assert.equal(comp.used, 200);
});

/** SystemUpdater **/
QUnit.test("SystemUpdater query builder", function(assert) {
	var su = new SystemUpdater(null);

	var keys = ['key1', 'key2'];
	assert.equal(su.updateQueryForComponents(keys), "?0=key1&1=key2");

	keys = ['with space'];
	assert.equal(su.updateQueryForComponents(keys), "?0=with%20space");

	keys = [];
	assert.equal(su.updateQueryForComponents(keys), "");
});

QUnit.asyncTest("SystemUpdater update components with 404", function(assert) {
	expect(2);

	$.mockjax({
		url : '*',
		status : 404,
		statusText : "Not Found",
		responseText : "Hello World!"
	});

	var sys = new System('Unit Test', 'http://example.org/');
	var su = new SystemUpdater(sys);
	su.updateComponents([], function(error) {
		assert.ok(error, "No error defined");

		pattern = /\(404\) Not Found/;
		pattern.exec(pattern);
		assert.ok(pattern.exec(error),
			"Correct HTTP status not found in '" + error + "'");

		QUnit.start();
	});

	$.mockjaxClear();
})

QUnit.asyncTest("SystemUpdater update components with success", function(assert) {
	expect(3);

	$.mockjax({
		url : '*',
		responseText : JSON.stringify({
			memory : {
				'total' : 500,
				'used' : 100
			}
		})
	});

	var sys = new System('Unit Test', 'http://example.org/');
	var mc = new MemoryComponent();
	sys.components[mc.key] = new MemoryComponent();

	var su = new SystemUpdater(sys);
	su.updateComponents([], function(error) {
		assert.ok(!error, "An error occurred");

		assert.equal(sys.components[mc.key].total, 500);
		assert.equal(sys.components[mc.key].used, 100);

		QUnit.start();
	});

	$.mockjaxClear();
})