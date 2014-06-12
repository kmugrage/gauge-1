Gem::Specification.new do |s|
    s.name        = "gauge-ruby"
    s.version     = '0.0.1'
    s.authors     = ["Gauge Team"]
    s.email       = ["gauge-team@thoughtworks.com"]

    s.summary     = "Ruby support for Gauge"
    s.description = "Adds Ruby support into Gauge tests"
    s.homepage    = "http://www.thoughtworks.com"
    s.files = Dir.glob("lib/**/*.rb")
    s.executables = Dir.glob("bin/**/*").map{ |f| File.basename(f)}
end

