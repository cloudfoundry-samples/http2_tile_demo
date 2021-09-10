#!/usr/bin/env ruby

url = "h2c.lightmediumorchid.cf-app.com"
iterations = 200
rate_limit = "1m"

puts "Doing #{iterations} iterations with a bandwidth limit of #{rate_limit} and 14kb of headers...\n\n"

start = Time.now
puts ">>> with http1"
iterations.times { system "curl -s --http1.1 -k -H @14k_of_headers.txt https://#{url}/ -o foo  --limit-rate #{rate_limit} > /dev/null" }
h1_average = (Time.now - start)/iterations
puts "Average Time: #{h1_average}\n\n"


start = Time.now
puts ">>> with http2"
iterations.times { system "curl -s -k -H @14k_of_headers.txt https://#{url}/ -o foo  --limit-rate #{rate_limit} > /dev/null" }
h2_average = (Time.now - start)/iterations
puts "Average Time: #{h2_average}"
puts "\n"
puts "#{(h1_average/h2_average).round(1)} times faster"
puts "#{(((h1_average/h2_average)-1)*100).round}% faster"
