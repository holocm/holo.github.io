#!/usr/bin/env perl

use strict;
use warnings;
use v5.32;

use JSON::PP qw(encode_json);
use Pod::Simple::SimpleTree;

# Pod::Simple::SimpleTree returns array-shaped nodes. Convert them into
# object-shaped nodes that are going to be easier to parse in Go.
sub transform_tree {
  my ($node) = @_;
  if (ref $node eq "ARRAY") {
    my ($name, $attrs, @subnodes) = @$node;
    $node = {
      name => $name,
      attrs => transform_attrs($attrs),
      children => [ map { transform_tree($_) } @subnodes ],
    };
  }
  return $node;
}

# Some attributes have values that are objects. This coerces all attributes
# into their string value.
sub transform_attrs {
  my ($attrs) = @_;
  return { map { $_ => "" . $attrs->{$_} } keys %$attrs };
}

my $parser = Pod::Simple::SimpleTree->new;
$parser->complain_stderr(1);
say(encode_json(transform_tree($parser->parse_file(shift)->root)));
