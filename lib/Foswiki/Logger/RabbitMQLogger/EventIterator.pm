package Foswiki::Logger::RabbitMQLogger::EventIterator;
use strict;
use warnings;

sub new {
    my ($class, $events, $version) = @_;

    unless($version) {
        $events = [map { _formatForOldApi($_) } @$events];
    }

    my $this = {
        events => $events,
        pos => 0,
    };

    return bless $this, $class;
}

# See Foswiki::Iterator::EventIterator::formatData(...)
sub _formatForOldApi {
    my $data = shift;

    my $extra = join(' ', @{$data->{extra} || []});

    return [
        $data->{epoch},
        $data->{user}       || '',
        $data->{action}     || '',
        $data->{webTopic}   || '',
        $extra,
        $data->{remoteAddr} || '',
        $data->{level},
    ];
}

sub hasNext {
    my ($this) = @_;

    return $this->{pos} < scalar @{$this->{events}};
}

sub next {
    my ($this) = @_;

    my $pos = $this->{pos}++;
    return $this->{events}->[$pos];
}

sub reset {
    my ($this) = @_;

    $this->{pos} = 0;
}

sub all {
    my ($this) = @_;

    return @{$this->{events}};
}

1;
