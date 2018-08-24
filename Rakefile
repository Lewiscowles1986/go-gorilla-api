# encoding: utf-8
require 'bundler'
require 'cucumber/rake/task'

Cucumber::Rake::Task.new

task :check do
  ENV['ROOT_URL'] = 'http://127.0.0.1:8080'
  if (ENV['WATCH'])
    Kernel.system("bundle exec guard -c")
  else
    Rake::Task['cucumber'].invoke
  end
end

task :help do
  puts """
Tasks
---
check - run tests
Environment Variables
---
WATCH - (guard) automatically rerun tests on file changes
FEATURE - path to specific cucumber feature to run
"""
end
