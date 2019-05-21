package Foswiki::Logger::RabbitMQLogger::Connection;

use strict;
use warnings;
use Net::AMQP::RabbitMQ;
use JSON;

use Exporter qw(import);
our @EXPORT = qw(createRabbitMQConnection);

sub DESTROY {
    my $this = shift;
    $this->{mq}->disconnect() if $this->{mq};
    delete $this->{mq};
}

my $routingKey = "qwiki.foswiki_perl";
my $exchange = 'ma_logs';


sub createRabbitMQConnection {
    return __PACKAGE__->new(@_);
}

sub new {
    my ($package, $fallback, $host, $user, $password) = @_;

    die "Please provide a fallback for RabbitMQ\n" unless $fallback;

    my $amqp = Net::AMQP::RabbitMQ->new();
    $amqp->connect($host, { user => $user, password => $password });
    $amqp->channel_open(1);

    $amqp->exchange_declare(
        1,
        $exchange,
        {
            exchange_type => 'topic',
            durable => 1,
            auto_delete => 0,
        },
    );

    my $this = {
        mq => $amqp,
        fallback => $fallback,
    };
    return bless $this, $package;
}

sub send {
    my ($this, $data) = @_;

    my $json = JSON->new();
    $json->allow_blessed(1);

    my $level =  $data->{log_data}->{level};
    eval {
        $this->{mq}->publish(
            1,
            "$routingKey.$level",
            $json->encode($data),
            {exchange => $exchange},
            {delivery_mode => 2},
        );
    };
    if($@) {
        $this->{fallback}->send({fields => ['Error while enqueueing message', $@]});
        $this->{fallback}->send($data);
    }
};

1;
