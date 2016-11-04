
#[]byte{0x5C},

def isprint(c)
  /[[:print:]]/ === c.chr
end

0.upto(255) {|c|
    if isprint(c)
        puts '[]byte{0x' + c.to_s(16) + '},'
    end
}
