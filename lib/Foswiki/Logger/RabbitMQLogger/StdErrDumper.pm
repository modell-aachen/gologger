package Foswiki::Logger::RabbitMQLogger::StdErrDumper;

use strict;
use warnings;
use JSON;

use Exporter qw(import);
our @EXPORT = qw(createSTDERRConnection);

sub createSTDERRConnection {
    return __PACKAGE__->new();
}

sub new {
    my $package = shift;

    my $this = {};
    return bless $this, $package;
}

sub send {
    my ($this, $data) = @_;

    my $time = $data->{time} || Foswiki::Time::formatTime(time(), 'iso', 'servertime');
    my $json = to_json($data);

    print STDERR "$time $json\n";
};

1;

