# 漢字
ary = []
while str = STDIN.gets
  code = str.split[5]
  if code =~ /^[0-9a-fA-F]+$/
      #bin = [code].pack("H*")
      #puts code + ' ' + bin
      ary << code
  end
end

ary.sort!

#ary.each {|v|
#    puts v
#}

# []byte{0xE3, 0x81, 0x82},

ary.each {|v|
    bary = []
    0.step(v.chomp.length - 2, 2) {|i|
        bary << '0x' + v[i,2]
    }
    puts '[]byte{' + bary.join(',') + '},'
}

