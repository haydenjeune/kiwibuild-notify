from diagrams import Cluster, Diagram, Edge
from diagrams.aws.compute import LambdaFunction
from diagrams.aws.database import DynamodbTable
from diagrams.aws.management import CloudwatchEventTimeBased
from diagrams.aws.integration import SimpleNotificationServiceSnsTopic, SimpleNotificationServiceSnsEmailNotification
from diagrams.custom import Custom

with Diagram("KiwiBuild Notifier", show=False):
    """
    Architectural diagram for the KiwiBuild Notifier system.
    """

    with Cluster("Internet"):
        internet = Custom("KiwiBuild Website", "docs/img/web.png")
    with Cluster("AWS"):
        trigger = CloudwatchEventTimeBased("Scheduled Trigger")
        scrape_fn = LambdaFunction("KiwiBuildScrapeFunction")
        table = DynamodbTable("Property")
        notify_fn = LambdaFunction("KiwiBuildNotifyFunction")
        topic = SimpleNotificationServiceSnsTopic("SNS Topic")
        notification = SimpleNotificationServiceSnsEmailNotification()

    trigger >> scrape_fn

    scrape_fn >> Edge(label="Store") >> table
    scrape_fn >> Edge(label="Read") >> internet

    table >> Edge(label="Stream trigger") >> notify_fn
    notify_fn >> Edge(label="Filter") >> topic
    topic >> notification
