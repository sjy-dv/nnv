# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: resourceCoordinator.proto
# Protobuf Python Version: 5.27.2
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    27,
    2,
    '',
    'resourceCoordinator.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x19resourceCoordinator.proto\x12\x15resourceCoordinatorV2\x1a\x1bgoogle/protobuf/empty.proto\"\xab\x02\n\nCollection\x12\x17\n\x0f\x63ollection_name\x18\x01 \x01(\t\x12\x31\n\x08\x64istance\x18\x02 \x01(\x0e\x32\x1f.resourceCoordinatorV2.Distance\x12\x39\n\x0cquantization\x18\x03 \x01(\x0e\x32#.resourceCoordinatorV2.Quantization\x12\x0b\n\x03\x64im\x18\x04 \x01(\r\x12\x14\n\x0c\x63onnectivity\x18\x05 \x01(\r\x12\x15\n\rexpansion_add\x18\x06 \x01(\r\x12\x18\n\x10\x65xpansion_search\x18\x07 \x01(\r\x12\r\n\x05multi\x18\x08 \x01(\x08\x12\x33\n\x07storage\x18\t \x01(\x0e\x32\".resourceCoordinatorV2.StorageType\"\x88\x01\n\x12\x43ollectionResponse\x12\x35\n\ncollection\x18\x01 \x01(\x0b\x32!.resourceCoordinatorV2.Collection\x12\x0e\n\x06status\x18\x02 \x01(\x08\x12+\n\x05\x65rror\x18\x03 \x01(\x0b\x32\x1c.resourceCoordinatorV2.Error\"\xba\x01\n\x10\x43ollectionDetail\x12\x35\n\ncollection\x18\x01 \x01(\x0b\x32!.resourceCoordinatorV2.Collection\x12\x17\n\x0f\x63ollection_size\x18\x02 \x01(\r\x12\x19\n\x11\x63ollection_memory\x18\x03 \x01(\x04\x12\x0e\n\x06status\x18\x04 \x01(\x08\x12+\n\x05\x65rror\x18\x05 \x01(\x0b\x32\x1c.resourceCoordinatorV2.Error\"{\n\x0e\x43ollectionList\x12\x35\n\ncollection\x18\x01 \x01(\x0b\x32!.resourceCoordinatorV2.Collection\x12\x17\n\x0f\x63ollection_size\x18\x02 \x01(\r\x12\x19\n\x11\x63ollection_memory\x18\x03 \x01(\x04\"\x8a\x01\n\x0f\x43ollectionLists\x12:\n\x0b\x63ollections\x18\x01 \x03(\x0b\x32%.resourceCoordinatorV2.CollectionList\x12\x0e\n\x06status\x18\x02 \x01(\x08\x12+\n\x05\x65rror\x18\x03 \x01(\x0b\x32\x1c.resourceCoordinatorV2.Error\"W\n\x18\x44\x65leteCollectionResponse\x12\x0e\n\x06status\x18\x01 \x01(\x08\x12+\n\x05\x65rror\x18\x02 \x01(\x0b\x32\x1c.resourceCoordinatorV2.Error\"<\n\x0e\x43ollectionName\x12\x17\n\x0f\x63ollection_name\x18\x01 \x01(\t\x12\x11\n\twith_size\x18\x02 \x01(\x08\"#\n\x0eGetCollections\x12\x11\n\twith_size\x18\x01 \x01(\x08\"G\n\x08Response\x12\x0e\n\x06status\x18\x01 \x01(\x08\x12+\n\x05\x65rror\x18\x02 \x01(\x0b\x32\x1c.resourceCoordinatorV2.Error\"T\n\x05\x45rror\x12\x15\n\rerror_message\x18\x01 \x01(\t\x12\x34\n\nerror_code\x18\x02 \x01(\x0e\x32 .resourceCoordinatorV2.ErrorCode\"\xbe\x01\n\nSystemInfo\x12\x0e\n\x06uptime\x18\x01 \x01(\x04\x12\x11\n\tcpu_load1\x18\x02 \x01(\x01\x12\x11\n\tcpu_load5\x18\x03 \x01(\x01\x12\x12\n\ncpu_load15\x18\x04 \x01(\x01\x12\x11\n\tmem_total\x18\x05 \x01(\x04\x12\x15\n\rmem_available\x18\x06 \x01(\x04\x12\x10\n\x08mem_used\x18\x07 \x01(\x04\x12\x10\n\x08mem_free\x18\x08 \x01(\x04\x12\x18\n\x10mem_used_percent\x18\t \x01(\x01*4\n\x0bStorageType\x12\x14\n\x10highspeed_memory\x10\x00\x12\x0f\n\x0bstable_disk\x10\x01*}\n\x08\x44istance\x12\x08\n\x04L2sq\x10\x00\x12\x06\n\x02Ip\x10\x01\x12\n\n\x06\x43osine\x10\x02\x12\r\n\tHaversine\x10\x03\x12\x0e\n\nDivergence\x10\x04\x12\x0b\n\x07Pearson\x10\x05\x12\x0b\n\x07Hamming\x10\x06\x12\x0c\n\x08Tanimoto\x10\x07\x12\x0c\n\x08Sorensen\x10\x08*M\n\x0cQuantization\x12\x08\n\x04None\x10\x00\x12\x08\n\x04\x42\x46\x31\x36\x10\x01\x12\x07\n\x03\x46\x31\x36\x10\x02\x12\x07\n\x03\x46\x33\x32\x10\x03\x12\x07\n\x03\x46\x36\x34\x10\x04\x12\x06\n\x02I8\x10\x05\x12\x06\n\x02\x42\x31\x10\x06*\x97\x01\n\tErrorCode\x12\r\n\tUNDEFINED\x10\x00\x12\r\n\tRPC_ERROR\x10\x01\x12!\n\x1d\x43OMMUNICATION_SHARD_RPC_ERROR\x10\x02\x12\x1d\n\x19\x43OMMUNICATION_SHARD_ERROR\x10\x03\x12\x11\n\rMARSHAL_ERROR\x10\x04\x12\x17\n\x13INTERNAL_FUNC_ERROR\x10\x05\x32\xfe\x05\n\x13ResourceCoordinator\x12\x38\n\x04Ping\x12\x16.google.protobuf.Empty\x1a\x16.google.protobuf.Empty\"\x00\x12\x62\n\x10\x43reateCollection\x12!.resourceCoordinatorV2.Collection\x1a).resourceCoordinatorV2.CollectionResponse\"\x00\x12l\n\x10\x44\x65leteCollection\x12%.resourceCoordinatorV2.CollectionName\x1a/.resourceCoordinatorV2.DeleteCollectionResponse\"\x00\x12\x61\n\rGetCollection\x12%.resourceCoordinatorV2.CollectionName\x1a\'.resourceCoordinatorV2.CollectionDetail\"\x00\x12\x64\n\x11GetAllCollections\x12%.resourceCoordinatorV2.GetCollections\x1a&.resourceCoordinatorV2.CollectionLists\"\x00\x12\x62\n\x0eLoadCollection\x12%.resourceCoordinatorV2.CollectionName\x1a\'.resourceCoordinatorV2.CollectionDetail\"\x00\x12]\n\x11ReleaseCollection\x12%.resourceCoordinatorV2.CollectionName\x1a\x1f.resourceCoordinatorV2.Response\"\x00\x12O\n\x10LoadResourceInfo\x12\x16.google.protobuf.Empty\x1a!.resourceCoordinatorV2.SystemInfo\"\x00\x42\x19Z\x17./resourceCoordinatorV2b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'resourceCoordinator_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\027./resourceCoordinatorV2'
  _globals['_STORAGETYPE']._serialized_start=1517
  _globals['_STORAGETYPE']._serialized_end=1569
  _globals['_DISTANCE']._serialized_start=1571
  _globals['_DISTANCE']._serialized_end=1696
  _globals['_QUANTIZATION']._serialized_start=1698
  _globals['_QUANTIZATION']._serialized_end=1775
  _globals['_ERRORCODE']._serialized_start=1778
  _globals['_ERRORCODE']._serialized_end=1929
  _globals['_COLLECTION']._serialized_start=82
  _globals['_COLLECTION']._serialized_end=381
  _globals['_COLLECTIONRESPONSE']._serialized_start=384
  _globals['_COLLECTIONRESPONSE']._serialized_end=520
  _globals['_COLLECTIONDETAIL']._serialized_start=523
  _globals['_COLLECTIONDETAIL']._serialized_end=709
  _globals['_COLLECTIONLIST']._serialized_start=711
  _globals['_COLLECTIONLIST']._serialized_end=834
  _globals['_COLLECTIONLISTS']._serialized_start=837
  _globals['_COLLECTIONLISTS']._serialized_end=975
  _globals['_DELETECOLLECTIONRESPONSE']._serialized_start=977
  _globals['_DELETECOLLECTIONRESPONSE']._serialized_end=1064
  _globals['_COLLECTIONNAME']._serialized_start=1066
  _globals['_COLLECTIONNAME']._serialized_end=1126
  _globals['_GETCOLLECTIONS']._serialized_start=1128
  _globals['_GETCOLLECTIONS']._serialized_end=1163
  _globals['_RESPONSE']._serialized_start=1165
  _globals['_RESPONSE']._serialized_end=1236
  _globals['_ERROR']._serialized_start=1238
  _globals['_ERROR']._serialized_end=1322
  _globals['_SYSTEMINFO']._serialized_start=1325
  _globals['_SYSTEMINFO']._serialized_end=1515
  _globals['_RESOURCECOORDINATOR']._serialized_start=1932
  _globals['_RESOURCECOORDINATOR']._serialized_end=2698
# @@protoc_insertion_point(module_scope)
