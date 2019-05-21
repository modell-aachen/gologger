package Foswiki::Logger::RabbitMQLogger::Reader;

use strict;
use warnings;
use Net::AMQP::RabbitMQ;
use JSON;
use Foswiki::Time;

sub read {
    my ($source, $startTime, $endTime, $levels, $host, $user, $password) = @_;

    $levels = [$levels] unless ref $levels;

    my $amqp = Net::AMQP::RabbitMQ->new();
    $amqp->connect($host, { user => $user, password => $password });
    $amqp->channel_open(1);

    my $rpc = $amqp->queue_declare(
        1,
        'rpc_queue',
        {
            durable => 1,
            exclusive => 0,
            auto_delete => 0,
        },
    );

    my $responseQueue = $amqp->queue_declare(
        1,
        '',
        {
            exclusive => 1,
            auto_delete => 1,
        },
    );
    $amqp->consume(1, $responseQueue);
    my $correlation_id = rand();

    my $formattedStartTime = Foswiki::Time::formatTime( $startTime, 'iso', 'servertime' );
    my $formattedEndTime = Foswiki::Time::formatTime( $endTime, 'iso', 'servertime' ) if $endTime;
    $amqp->publish(
        1,
        'rpc_queue',
        to_json(
            {
                levels => $levels,
                start_time => $formattedStartTime,
                end_time => $formattedEndTime,
                source => $source,
            }
        ),
        undef,
        {
            reply_to => $responseQueue,
            correlation_id => $correlation_id,
            delivery_mode => 2,
        },
    );

    my $received;
    #print STDERR "waiting fr ressponse on $responseQueue with id $correlation_id\n";
    do {
        $received = $amqp->recv(0);
    } while ($received->{props}{correlation_id} ne $correlation_id);

    #print STDERR "disconnecting\n";
    $amqp->disconnect();

    my $response = from_json($received->{body});

    my @events = map {
        my $event = $_;
        my $data = $event->{Misc} || {};
        $data->{level} = $event->{Level};
        $data->{epoch} = Foswiki::Time::parseTime($event->{Time});
        $data->{extra} = $event->{Fields} if $event->{Fields};
        $data;
    } @$response;
    return \@events;
}

1;

