import grpc
from concurrent import futures
import time
import questgen_pb2 as pb2
import questgen_pb2_grpc as pb2_grpc



class QuestGenerator(pb2_grpc.QuestGenServiceServicer):

    def __init__(self, *args, **kwargs):
        pass

    def QuestGen(self, request, context):

        # get the string from the incoming request
        requestData = request.request
        result = f'Hello I am up and running received "{requestData}" message from you'
        result = {'response': result}

        return pb2.QuestGenResponse(**result)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    pb2_grpc.add_QuestGenServiceServicer_to_server(QuestGenerator(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    serve()